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
	CPP_CODE_COUNT_STATE_INIT               = 0
	CPP_CODE_COUNT_STATE_SLASH              = 1
	CPP_CODE_COUNT_STATE_BLOCK_COMMENT      = 2
	CPP_CODE_COUNT_STATE_BLOCK_COMMENT_STAR = 3
	CPP_CODE_COUNT_STATE_LINE_COMMENT       = 4
	CPP_CODE_COUNT_STATE_CODE               = 5
	CPP_CODE_COUNT_STATE_STRING             = 6
	CPP_CODE_COUNT_STATE_STRING_ESCAPE      = 7
)

type CppCodeCounter struct {
	state int
}

func NewCppCodeCounter() *CppCodeCounter {
	return &CppCodeCounter{}
}

func (c *CppCodeCounter) Clear() { c.state = CPP_CODE_COUNT_STATE_INIT }

func (c *CppCodeCounter) ParseFile(filename string) (stat CodeStat, ok bool) {
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

func (c *CppCodeCounter) ParseLine(line string) (stat CodeStat) {
	stat.Total = 1
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		stat.Blank = 1
		return stat
	}

	hasCode := false
	hasComment := false

	if c.state != CPP_CODE_COUNT_STATE_BLOCK_COMMENT && c.state != CPP_CODE_COUNT_STATE_STRING && c.state != CPP_CODE_COUNT_STATE_STRING_ESCAPE {
		c.state = CPP_CODE_COUNT_STATE_INIT
	}

	for _, v := range line {
		if c.state == CPP_CODE_COUNT_STATE_LINE_COMMENT {
			break
		}

		if hasCode && hasComment {
			break
		}

		switch c.state {
		case CPP_CODE_COUNT_STATE_INIT:
			switch v {
			case ' ':
				break
			case '\t':
				break
			case '/':
				c.state = CPP_CODE_COUNT_STATE_SLASH
			case '"':
				c.state = CPP_CODE_COUNT_STATE_STRING
				hasCode = true
			default:
				c.state = CPP_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case CPP_CODE_COUNT_STATE_SLASH:
			switch v {
			case '*':
				c.state = CPP_CODE_COUNT_STATE_BLOCK_COMMENT
				hasComment = true
			case '/':
				c.state = CPP_CODE_COUNT_STATE_LINE_COMMENT
				hasComment = true
			case '"':
				c.state = CPP_CODE_COUNT_STATE_STRING
				hasCode = true
			default:
				c.state = CPP_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case CPP_CODE_COUNT_STATE_BLOCK_COMMENT:
			hasComment = true
			if v == '*' {
				c.state = CPP_CODE_COUNT_STATE_BLOCK_COMMENT_STAR
			}

		case CPP_CODE_COUNT_STATE_BLOCK_COMMENT_STAR:
			if v == '/' {
				c.state = CPP_CODE_COUNT_STATE_INIT
			}

		case CPP_CODE_COUNT_STATE_CODE:
			switch v {
			case ' ':
				fallthrough
			case '\t':
				c.state = CPP_CODE_COUNT_STATE_INIT
			case '/':
				c.state = CPP_CODE_COUNT_STATE_SLASH
			case '"':
				c.state = CPP_CODE_COUNT_STATE_STRING
				hasCode = true
			}

		case CPP_CODE_COUNT_STATE_STRING:
			hasCode = true
			switch v {
			case '\\':
				c.state = CPP_CODE_COUNT_STATE_STRING_ESCAPE
			case '"':
				c.state = CPP_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case CPP_CODE_COUNT_STATE_STRING_ESCAPE:
			c.state = CPP_CODE_COUNT_STATE_STRING

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
