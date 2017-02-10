package counter

import (
	"fmt"
)

type CodeStat struct {
	Total   int
	Code    int
	Comment int
	Blank   int
}

func (c *CodeStat) String() string {
	return fmt.Sprintf("Total = %6d, Code = %6d, Comment = %6d, Blank = %6d, CommentPercent = %2.2f%%",
		c.Total, c.Code, c.Comment, c.Blank, c.CommentPercent())
}

func (c *CodeStat) Add(rhs *CodeStat) {
	c.Total += rhs.Total
	c.Code += rhs.Code
	c.Comment += rhs.Comment
	c.Blank += rhs.Blank
}

func (c *CodeStat) CommentPercent() float64 {
	if (c.Code + c.Comment) == 0 {
		return 0.0
	}
	return float64(c.Comment) / float64(c.Code+c.Comment) * 100
}

type CodeCounter interface {
	Clear()
	ParseFile(filename string) (stat CodeStat, ok bool)
	ParseLine(line string) (stat CodeStat)
}
