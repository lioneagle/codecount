package counter

import (
	//"fmt"
	"os"
	//"path/filepath"
	"testing"
)

func TestErlangCodeCounterParseLine(t *testing.T) {
	testdata := []struct {
		line string
		stat CodeStat
	}{
		{" \t", CodeStat{Total: 1, Blank: 1}},
		{"\t abc \t", CodeStat{Total: 1, Code: 1}},
		{"ab%", CodeStat{Total: 1, Code: 1, Comment: 1}},
		{"%ab%", CodeStat{Total: 1, Comment: 1}},
		{" \"%\"\tcc\t ff\"\\t\" d", CodeStat{Total: 1, Code: 1, Comment: 0}},
		{" '%'\tcc\t ff\"\\t\" d", CodeStat{Total: 1, Code: 1, Comment: 0}},
		{" aa'%'\tcc\t ff\"\\t\" d", CodeStat{Total: 1, Code: 1, Comment: 0}},
	}

	for i, v := range testdata {
		counter, _ := NewCodeCounterFactory().NewCounter("erlang")

		stat := counter.ParseLine(v.line)

		if stat != v.stat {
			t.Errorf("TestErlangCodeCounterParseLine[%d] failed, stat = %s, wanted = %s", i, stat.String(), v.stat.String())
			continue
		}
	}
}

func TestErlangCodeCounterParseFile(t *testing.T) {
	filename := os.Args[len(os.Args)-1] + "\\src\\testdata\\test1.erl"

	counter, _ := NewCodeCounterFactory().NewCounter("erlang")
	wanted := CodeStat{Total: 11, Code: 5, Comment: 4, Blank: 3}

	stat, ok := ParseFile(counter, filename)
	if !ok {
		t.Errorf("TestErlangCodeCounterParseFile failed, ParseFile failed")
		return
	}

	if stat != wanted {
		t.Errorf("TestErlangCodeCounterParseFile failed, stat = %s, wanted = %s", stat.String(), wanted.String())
		return
	}
}
