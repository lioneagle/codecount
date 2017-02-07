package counter

import (
	//"fmt"
	//"os"
	//"path/filepath"
	"testing"
)

func TestCppCodeCounterParseLine(t *testing.T) {
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
		{" ff/\"as/\"", CodeStat{Total: 1, Code: 1}},
		{"var x= 3", CodeStat{Total: 1, Code: 1}},
		{"a\"\\\"", CodeStat{Total: 1, Code: 1}},
	}

	for i, v := range testdata {
		counter, _ := NewCodeCounterFactory().NewCounter("cpp")

		stat := counter.ParseLine(v.line)

		if stat != v.stat {
			t.Errorf("TestCppCodeCounterParseLine[%d] failed, stat = %s, wanted = %s", i, stat.String(), v.stat.String())
			continue
		}
	}
}

func TestCppCodeCounterParseFile(t *testing.T) {
	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//fmt.Printf("arg[0] = %s, dir = %s\n", os.Args[0], dir)

	//filename := "testdata/test1.go"
	filename := "F:/dev/go_code/src/codecount/src/testdata/test1.cpp"

	counter, _ := NewCodeCounterFactory().NewCounter("cpp")
	wanted := CodeStat{Total: 12, Code: 9, Comment: 3, Blank: 1}

	stat, ok := counter.ParseFile(filename)
	if !ok {
		t.Errorf("TestCppCodeCounterParseFile failed, ParseFile failed")
		return
	}

	if stat != wanted {
		t.Errorf("TestCppCodeCounterParseFile failed, stat = %s, wanted = %s", stat.String(), wanted.String())
		return
	}
}
