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

func main() {
	filter := "*.cpp;*.cxx;*.go"
	filters := strings.Split(filter, ";")

	codetypeMap := NewCodeTypeMap()
	codetypeMap.AddCodeType("*.cpp;*.cxx;*.hpp;*.hxx;*.c++;*.cc;*.c;*.h", "cpp")
	codetypeMap.AddCodeType("*.go", "go")

	root := "F:/dev/go_code/src/codecount"

	files := &fileList{}

	getFiles(root, filters, files)

	total_stat := &counter.CodeStat{}

	factory := counter.NewCodeCounterFactory()

	for _, v := range files.data {
		codetype, ok := codetypeMap.maps[strings.ToLower(filepath.Ext(v.filename)[1:])]
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
		total_stat.Add(&stat)

		fmt.Printf("%s: %s\n", v.filename, stat.String())
	}

	fmt.Printf("total %d files: %s\n", len(files.data), total_stat.String())
}
