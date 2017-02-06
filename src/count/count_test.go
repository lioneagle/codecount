package count

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	testdata := []struct {
		line string
		stat CodeStat
	}{
		{" \t", CodeStat{Total: 1, Blank: 1}},
		{"ab/c", CodeStat{Total: 1, Code: 1}},
		{"ab//", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{"//ab//", CodeStat{Total: 1, Comment: 1}},
		{"ab/*", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{"/*ab", CodeStat{Total: 1, Comment: 1}},
		{"ab/*tt**/ cc ", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{" /*aa**/ \tcc\t ff d", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{" \"//\"/*aa*/ \tcc\t ff\"\\t\" d", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{" '//\\a\"'/*aa*/ \tcc\t ff", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{" ff`as/`", CodeStat{Total: 1, Code: 1}},
		{" ff/`as/`", CodeStat{Total: 1, Code: 1}},
		{" ff/\"as/\"", CodeStat{Total: 1, Code: 1}},
		{"var x= 3", CodeStat{Total: 1, Code: 1}},
		{"`\\t`", CodeStat{Total: 1, Code: 1}},
		{"a\"\\\"", CodeStat{Total: 1, Code: 1}},
	}

	for i, v := range testdata {
		counter := &GoCounter{}

		stat := counter.ParseLine(v.line)

		if stat != v.stat {
			t.Errorf("TestParseLine[%d] failed, stat = %s, wanted = %s", i, stat.String(), v.stat.String())
			continue
		}
	}

}
