package counter

import (
	//"fmt"
	"strings"
)

const (
	ERLANG_CODE_COUNT_STATE_INIT         = 0
	ERLANG_CODE_COUNT_STATE_LINE_COMMENT = 1
	ERLANG_CODE_COUNT_STATE_CODE         = 2
	ERLANG_CODE_COUNT_STATE_STRING       = 3
	ERLANG_CODE_COUNT_STATE_ATOM         = 4
)

type ErlangCodeCounter struct {
	state int
}

func NewErlangCodeCounter() *ErlangCodeCounter {
	return &ErlangCodeCounter{}
}

func (c *ErlangCodeCounter) Clear() { c.state = ERLANG_CODE_COUNT_STATE_INIT }

func (c *ErlangCodeCounter) ParseLine(line string) (stat CodeStat) {
	stat.Total = 1
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		stat.Blank = 1
		return stat
	}

	hasCode := false
	hasComment := false

	/*if c.state != ERLANG_CODE_COUNT_STATE_STRING && c.state != ERLANG_CODE_COUNT_STATE_ATOM {
		c.state = ERLANG_CODE_COUNT_STATE_INIT
	}*/

	for _, v := range line {
		//fmt.Printf("v = %c, state = %d\n", v, c.state)

		if c.state == ERLANG_CODE_COUNT_STATE_LINE_COMMENT {
			break
		}

		switch c.state {
		case ERLANG_CODE_COUNT_STATE_INIT:
			switch v {
			case '%':
				c.state = ERLANG_CODE_COUNT_STATE_LINE_COMMENT
				hasComment = true
			case '"':
				c.state = ERLANG_CODE_COUNT_STATE_STRING
				hasCode = true
			case '\'':
				c.state = ERLANG_CODE_COUNT_STATE_ATOM
				hasCode = true
			default:
				c.state = ERLANG_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case ERLANG_CODE_COUNT_STATE_CODE:
			switch v {
			case '%':
				c.state = ERLANG_CODE_COUNT_STATE_LINE_COMMENT
				hasComment = true
			case '"':
				c.state = ERLANG_CODE_COUNT_STATE_STRING
				hasCode = true
			case '\'':
				c.state = ERLANG_CODE_COUNT_STATE_ATOM
				hasCode = true
			default:
				c.state = ERLANG_CODE_COUNT_STATE_CODE
				hasCode = true
			}

		case ERLANG_CODE_COUNT_STATE_STRING:
			hasCode = true
			if v == '"' {
				c.state = ERLANG_CODE_COUNT_STATE_CODE
			}

		case ERLANG_CODE_COUNT_STATE_ATOM:
			hasCode = true
			if v == '\'' {
				c.state = ERLANG_CODE_COUNT_STATE_CODE
			}

		}
	}

	if hasCode {
		stat.Code = 1
	}

	if hasComment {
		stat.Comment = 1
	}

	switch c.state {
	case ERLANG_CODE_COUNT_STATE_STRING:
		break
	case ERLANG_CODE_COUNT_STATE_ATOM:
		break
	default:
		c.state = ERLANG_CODE_COUNT_STATE_INIT
	}

	return stat
}
