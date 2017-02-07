package counter

import (
	"bufio"
	//"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	GO_CODE_COUNT_STATE_INIT                = 0
	GO_CODE_COUNT_STATE_SLASH               = 1
	GO_CODE_COUNT_STATE_BLOCK_COMMENT       = 2
	GO_CODE_COUNT_STATE_BLOCK_COMMENT_STAR  = 3
	GO_CODE_COUNT_STATE_LINE_COMMENT        = 4
	GO_CODE_COUNT_STATE_CODE                = 5
	GO_CODE_COUNT_STATE_LINE_STRING         = 6
	GO_CODE_COUNT_STATE_LINE_STRING_ESCAPE  = 7
	GO_CODE_COUNT_STATE_BLOCK_STRING        = 8
	GO_CODE_COUNT_STATE_BLOCK_STRING_ESCAPE = 9
)

type GoCounter struct {
	state int
}

func NewGoCounter() *GoCounter {
	return &GoCounter{}
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

	if c.state != GO_CODE_COUNT_STATE_BLOCK_COMMENT && c.state != GO_CODE_COUNT_STATE_BLOCK_STRING {
		c.state = GO_CODE_COUNT_STATE_INIT
	}

	for _, v := range line {
		if c.state == GO_CODE_COUNT_STATE_LINE_COMMENT {
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
				c.state = GO_CODE_COUNT_STATE_LINE_STRING
				hasCode = true
			case '`':
				c.state = GO_CODE_COUNT_STATE_BLOCK_STRING
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_SLASH:
			switch v {
			case '*':
				c.state = GO_CODE_COUNT_STATE_BLOCK_COMMENT
				hasComment = true
			case '/':
				c.state = GO_CODE_COUNT_STATE_LINE_COMMENT
				hasComment = true
			case '"':
				c.state = GO_CODE_COUNT_STATE_LINE_STRING
				hasCode = true
			case '`':
				c.state = GO_CODE_COUNT_STATE_BLOCK_STRING
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_BLOCK_COMMENT:
			hasComment = true
			if v == '*' {
				c.state = GO_CODE_COUNT_STATE_BLOCK_COMMENT_STAR
			}

		case GO_CODE_COUNT_STATE_BLOCK_COMMENT_STAR:
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
				c.state = GO_CODE_COUNT_STATE_LINE_STRING
				hasCode = true
			case '`':
				c.state = GO_CODE_COUNT_STATE_BLOCK_STRING
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_LINE_STRING:
			switch v {
			case '\\':
				c.state = GO_CODE_COUNT_STATE_LINE_STRING_ESCAPE
			case '"':
				c.state = GO_CODE_COUNT_STATE_CODE
			default:
				c.state = GO_CODE_COUNT_STATE_LINE_STRING
			}

		case GO_CODE_COUNT_STATE_BLOCK_STRING:
			hasCode = true
			switch v {
			case '\\':
				c.state = GO_CODE_COUNT_STATE_BLOCK_STRING_ESCAPE
			case '`':
				c.state = GO_CODE_COUNT_STATE_CODE
				hasCode = true
			default:
				c.state = GO_CODE_COUNT_STATE_BLOCK_STRING
			}

		case GO_CODE_COUNT_STATE_LINE_STRING_ESCAPE:
			c.state = GO_CODE_COUNT_STATE_LINE_STRING

		case GO_CODE_COUNT_STATE_BLOCK_STRING_ESCAPE:
			c.state = GO_CODE_COUNT_STATE_BLOCK_STRING

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
