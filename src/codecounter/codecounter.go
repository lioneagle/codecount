package main

import (
	"counter"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fileInfo struct {
	filename string
	codetype string
	stat     counter.CodeStat
}

type fileList struct {
	data []*fileInfo
}

func getFiles(root string, filters []string, files *fileList) error {
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		for _, v := range filters {
			if ok, _ := filepath.Match(v, f.Name()); ok {
				fileinfo := &fileInfo{filename: path}
				files.data = append(files.data, fileinfo)
				//fmt.Println(fileinfo.filename)
			}

		}
		return nil
	}

	return filepath.Walk(root, walkFunc)
}

type CodeTypeMap struct {
	maps map[string]string
}

func NewCodeTypeMap() *CodeTypeMap {
	return &CodeTypeMap{maps: make(map[string]string)}
}

func (c *CodeTypeMap) AddCodeType(filters string, codetype string) {
	ext := strings.Split(filters, ";")

	for _, v := range ext {
		c.maps[strings.ToLower(filepath.Ext(v)[1:])] = strings.ToLower(codetype)
	}
}

type CodeTypeStat struct {
	filenum int
	stat    counter.CodeStat
}

type CodeTypeStats struct {
	maps map[string]*CodeTypeStat
}

func NewCodeTypeStats() *CodeTypeStats {
	return &CodeTypeStats{maps: make(map[string]*CodeTypeStat)}
}

func main() {
	filter := "*.cpp;*.cxx;*.go"
	filters := strings.Split(filter, ";")

	codetypeMap := NewCodeTypeMap()
	codetypeMap.AddCodeType("*.cpp;*.cxx;*.hpp;*.hxx;*.c++;*.cc;*.c;*.h", "cpp")
	codetypeMap.AddCodeType("*.go", "go")

	codetypes := []string{"cpp", "go"}

	root := "F:/dev/go_code/src/codecount"

	files := &fileList{}

	getFiles(root, filters, files)

	totlStat := &counter.CodeStat{}
	codetypeStats := NewCodeTypeStats()

	factory := counter.NewCodeCounterFactory()

	for _, v := range files.data {
		v.codetype = strings.ToLower(filepath.Ext(v.filename)[1:])
		codetype, ok := codetypeMap.maps[v.codetype]
		if !ok {
			log.Printf("ERROR: unknown code type for %s", v.filename)
			continue
		}

		c, ok := factory.NewCounter(codetype)
		if !ok {
			log.Printf("ERROR: cannot get codecounter for %s", v.filename)
			continue
		}

		stat, ok := c.ParseFile(v.filename)
		if !ok {
			log.Printf("ERROR: parse file %s failed", v.filename)
		}
		totlStat.Add(&stat)

		codetypeStat, ok := codetypeStats.maps[v.codetype]
		if !ok {
			codetypeStat = &CodeTypeStat{}
			codetypeStats.maps[v.codetype] = codetypeStat
		}

		codetypeStat.filenum++
		codetypeStat.stat.Add(&stat)

		fmt.Printf("%s: %s\n", v.filename, stat.String())
	}

	for _, v := range codetypes {
		fmt.Printf("total %d %s files: %s\n", codetypeStats.maps[v].filenum, v, codetypeStats.maps[v].stat.String())
	}
	fmt.Printf("total %d files: %s\n", len(files.data), totlStat.String())
}
