package count

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CodeStat struct {
	Total   int
	Code    int
	Comment int
	Blank   int
}

func (c *CodeStat) String() string {
	return fmt.Sprintf("Total = %d, Code = %d, Comment = %d, Blank = %d, CommentPercent = %.2f%",
		c.Total, c.Code, c.Comment, c.Blank, c.CommentPercent())
}

func (c *CodeStat) Add(rhs *CodeStat) {
	c.Total += rhs.Total
	c.Code += rhs.Code
	c.Comment += rhs.Comment
	c.Blank += rhs.Blank
}

func (c *CodeStat) CommentPercent() float64 {
	return float64(c.Comment) / float64(c.Code+c.Comment)
}

const (
	GO_CODE_COUNT_STATE_INIT           = 0
	GO_CODE_COUNT_STATE_SLASH          = 1
	GO_CODE_COUNT_STATE_COMMENT1       = 2
	GO_CODE_COUNT_STATE_COMMENT1_STAR  = 3
	GO_CODE_COUNT_STATE_COMMENT2       = 4
	GO_CODE_COUNT_STATE_CODE           = 5
	GO_CODE_COUNT_STATE_STRING1        = 6
	GO_CODE_COUNT_STATE_STRING1_ESCAPE = 7
	GO_CODE_COUNT_STATE_STRING2        = 8
	GO_CODE_COUNT_STATE_STRING2_ESCAPE = 9
)

type GoCounter struct {
	state int
}

func (c *GoCounter) Clear() { c.state = GO_CODE_COUNT_STATE_INIT }

func (c *GoCounter) ParseFile(filename string) (stat CodeStat, ok bool) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("ERROR: cannot open file %s", filename)
		return stat, false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		lineStat := c.ParseLine(line)
		stat.Add(&lineStat)
	}
	return stat, true
}

func (c *GoCounter) ParseLine(line string) (stat CodeStat) {
	stat.Total = 1
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		stat.Blank = 1
		return stat
	}

	hasCode := false
	hasComment := false

	if c.state != GO_CODE_COUNT_STATE_COMMENT1 {
		c.state = GO_CODE_COUNT_STATE_INIT
	}

	for _, v := range line {
		if c.state == GO_CODE_COUNT_STATE_COMMENT2 {
			break
		}

		if hasCode && hasComment {
			break
		}

		switch c.state {
		case GO_CODE_COUNT_STATE_INIT:
			switch v {
			case ' ':
				break
			case '\t':
				break
			case '/':
				c.state = GO_CODE_COUNT_STATE_SLASH
			case '"':
				c.state = GO_CODE_COUNT_STATE_STRING1
				hasCode = true
			case '`':
				c.state = GO_CODE_COUNT_STATE_STRING2
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_SLASH:
			switch v {
			case '*':
				c.state = GO_CODE_COUNT_STATE_COMMENT1
				hasComment = true
			case '/':
				c.state = GO_CODE_COUNT_STATE_COMMENT2
				hasComment = true
			case '"':
				c.state = GO_CODE_COUNT_STATE_STRING1
				hasCode = true
			case '`':
				c.state = GO_CODE_COUNT_STATE_STRING2
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_COMMENT1:
			if v == '*' {
				c.state = GO_CODE_COUNT_STATE_COMMENT1_STAR
			}

		case GO_CODE_COUNT_STATE_COMMENT1_STAR:
			if v == '/' {
				c.state = GO_CODE_COUNT_STATE_INIT
			}

		case GO_CODE_COUNT_STATE_CODE:
			switch v {
			case ' ':
				fallthrough
			case '\t':
				c.state = GO_CODE_COUNT_STATE_INIT
			case '/':
				c.state = GO_CODE_COUNT_STATE_SLASH
			case '"':
				c.state = GO_CODE_COUNT_STATE_STRING1
				hasCode = true
			case '`':
				c.state = GO_CODE_COUNT_STATE_STRING2
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_STRING1:
			switch v {
			case '\\':
				c.state = GO_CODE_COUNT_STATE_STRING1_ESCAPE
			case '"':
				c.state = GO_CODE_COUNT_STATE_CODE
			default:
				c.state = GO_CODE_COUNT_STATE_STRING1
			}

		case GO_CODE_COUNT_STATE_STRING2:
			switch v {
			case '\\':
				c.state = GO_CODE_COUNT_STATE_STRING2_ESCAPE
			case '`':
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_STRING2
			}

		case GO_CODE_COUNT_STATE_STRING1_ESCAPE:
			c.state = GO_CODE_COUNT_STATE_STRING1

		case GO_CODE_COUNT_STATE_STRING2_ESCAPE:
			c.state = GO_CODE_COUNT_STATE_STRING2

		}
	}

	if hasCode {
		stat.Code = 1
	}

	if hasComment {
		stat.Comment = 1
	}

	return stat
}
