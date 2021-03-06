package counter

import (
	//"fmt"
	"strings"
)

const (
	GO_CODE_COUNT_STATE_INIT               = 0
	GO_CODE_COUNT_STATE_SLASH              = 1
	GO_CODE_COUNT_STATE_BLOCK_COMMENT      = 2
	GO_CODE_COUNT_STATE_BLOCK_COMMENT_STAR = 3
	GO_CODE_COUNT_STATE_LINE_COMMENT       = 4
	GO_CODE_COUNT_STATE_CODE               = 5
	GO_CODE_COUNT_STATE_LINE_STRING        = 6
	GO_CODE_COUNT_STATE_LINE_STRING_ESCAPE = 7
	GO_CODE_COUNT_STATE_BLOCK_STRING       = 8
	GO_CODE_COUNT_STATE_CHAR               = 9
	GO_CODE_COUNT_STATE_CHAR_ESCAPE        = 10
)

type GoCodeCounter struct {
	state int
}

func NewGoCodeCounter() *GoCodeCounter {
	return &GoCodeCounter{}
}

func (c *GoCodeCounter) Clear() { c.state = GO_CODE_COUNT_STATE_INIT }

func (c *GoCodeCounter) ParseLine(line string) (stat CodeStat) {
	stat.Total = 1
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		switch c.state {
		case GO_CODE_COUNT_STATE_BLOCK_COMMENT:
			stat.Comment = 1
		case GO_CODE_COUNT_STATE_BLOCK_STRING:
			stat.Code = 1
		default:
			stat.Blank = 1
		}
		return stat
	}

	hasCode := false
	hasComment := false

	/*if c.state != GO_CODE_COUNT_STATE_BLOCK_COMMENT && c.state != GO_CODE_COUNT_STATE_BLOCK_STRING {
		c.state = GO_CODE_COUNT_STATE_INIT
	}*/

	for _, v := range line {
		if c.state == GO_CODE_COUNT_STATE_LINE_COMMENT {
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
			case '\'':
				c.state = GO_CODE_COUNT_STATE_CHAR
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
			case '\'':
				c.state = GO_CODE_COUNT_STATE_CHAR
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
			} else if v != '*' {
				c.state = GO_CODE_COUNT_STATE_BLOCK_COMMENT
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
			case '\'':
				c.state = GO_CODE_COUNT_STATE_CHAR
				hasCode = true
			}

		case GO_CODE_COUNT_STATE_LINE_STRING:
			switch v {
			case '\\':
				c.state = GO_CODE_COUNT_STATE_LINE_STRING_ESCAPE
			case '"':
				c.state = GO_CODE_COUNT_STATE_CODE
			}

		case GO_CODE_COUNT_STATE_LINE_STRING_ESCAPE:
			c.state = GO_CODE_COUNT_STATE_LINE_STRING

		case GO_CODE_COUNT_STATE_BLOCK_STRING:
			hasCode = true
			if v == '`' {
				c.state = GO_CODE_COUNT_STATE_CODE
			}

		case GO_CODE_COUNT_STATE_CHAR:
			switch v {
			case '\\':
				c.state = GO_CODE_COUNT_STATE_CHAR_ESCAPE
			case '\'':
				c.state = GO_CODE_COUNT_STATE_CODE
			}

		case GO_CODE_COUNT_STATE_CHAR_ESCAPE:
			c.state = GO_CODE_COUNT_STATE_CHAR
		}
	}

	if hasCode {
		stat.Code = 1
	}

	if hasComment {
		stat.Comment = 1
	}

	switch c.state {
	case GO_CODE_COUNT_STATE_BLOCK_COMMENT:
		break
	case GO_CODE_COUNT_STATE_BLOCK_COMMENT_STAR:
		c.state = GO_CODE_COUNT_STATE_BLOCK_COMMENT
	case GO_CODE_COUNT_STATE_BLOCK_STRING:
		break
	default:
		c.state = GO_CODE_COUNT_STATE_INIT
	}

	return stat
}
