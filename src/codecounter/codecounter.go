package main

import (
	"counter"
	"flag"
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
		if v == "" {
			continue
		}
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

	codeConfig := []struct {
		name          string
		extConfig     string
		extConfigDesc string
	}{
		{"cpp", "*.cpp;*.cxx;*.hpp;*.hxx;*.c++;*.cc", "extions for c/c++ files"},
		{"c", "*.c;*.h", "extions for c files"},
		{"go", "*.go", "extions for go files"},
	}

	root := flag.String("path", ".", "path for code")
	filter := flag.String("filter", "*.cpp;*.cxx;*.hpp;*.hxx;*.c++;*.cc;*.c;*.h;*.go", "file filters")
	showEachFile := flag.Bool("show", false, "disable show each file stat")

	exts := make([]*string, 0)
	for _, v := range codeConfig {
		exts = append(exts, flag.String(v.name, v.extConfig, v.extConfigDesc))
	}

	flag.Parse()

	_, err := os.Stat(*root)
	if os.IsNotExist(err) {
		fmt.Printf("ERROR: path \"%s\" is not exist", *root)
		return
	}

	codetypeMap := NewCodeTypeMap()
	codetypeStats := NewCodeTypeStats()
	for i, v := range codeConfig {
		codetypeMap.AddCodeType(*exts[i], v.name)
		codetypeStats.maps[v.name] = &CodeTypeStat{}
	}

	factory := counter.NewCodeCounterFactory()

	files := &fileList{}
	getFiles(*root, strings.Split(*filter, ";"), files)

	totlStat := &counter.CodeStat{}

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

		codetypeStat, ok := codetypeStats.maps[codetype]
		if !ok {
			codetypeStat = &CodeTypeStat{}
			codetypeStats.maps[codetype] = codetypeStat
		}

		codetypeStat.filenum++
		codetypeStat.stat.Add(&stat)

		if *showEachFile {
			fmt.Printf("%s: %s\n", v.filename, stat.String())
		}
	}

	for _, v := range codeConfig {
		if codetypeStats.maps[v.name].filenum > 0 {
			fmt.Printf("total %d %s files: %s\n", codetypeStats.maps[v.name].filenum, v.name, codetypeStats.maps[v.name].stat.String())
		}
	}
	fmt.Printf("total %d files: %s\n", len(files.data), totlStat.String())
}
