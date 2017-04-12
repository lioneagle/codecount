package counter

import (
	"bufio"
	"fmt"
	"strconv"
	//"fmt"
	"io"
	"log"
	"os"
)

type CodeStat struct {
	Total   int
	Code    int
	Comment int
	Blank   int
}

func (codeStat *CodeStat) String() string {
	return fmt.Sprintf("Total = %6d, Code = %6d, Comment = %6d, Blank = %6d, CommentPercent = %2.2f%%",
		codeStat.Total, codeStat.Code, codeStat.Comment, codeStat.Blank, codeStat.CommentPercent())
}

func (codeStat *CodeStat) StringSlice() []string {
	return []string{
		strconv.Itoa(codeStat.Total),
		strconv.Itoa(codeStat.Code),
		strconv.Itoa(codeStat.Comment),
		strconv.Itoa(codeStat.Blank),
		fmt.Sprintf("%2.2f%%", codeStat.CommentPercent()),
	}

}

func (codeStat *CodeStat) Add(rhs *CodeStat) {
	codeStat.Total += rhs.Total
	codeStat.Code += rhs.Code
	codeStat.Comment += rhs.Comment
	codeStat.Blank += rhs.Blank
}

func (codeStat *CodeStat) CommentPercent() float64 {
	if (codeStat.Code + codeStat.Comment) == 0 {
		return 0.0
	}
	return float64(codeStat.Comment) / float64(codeStat.Code+codeStat.Comment) * 100
}

type CodeCounter interface {
	Clear()
	ParseLine(line string) (stat CodeStat)
}

func ParseFile(counter CodeCounter, filename string) (stat CodeStat, ok bool) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("ERROR: cannot open file %s", filename)
		return stat, false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && io.EOF != err {
			break
		}
		lineStat := counter.ParseLine(line)
		stat.Add(&lineStat)
		//fmt.Printf("line = %s\nlineStat = %s, state = %d\n", strings.TrimSpace(line), lineStat.String(), c.state)

		if io.EOF == err {
			break
		}
	}
	return stat, true
}
