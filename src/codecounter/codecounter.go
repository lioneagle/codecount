package main

import (
	"counter"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type fileInfo struct {
	filename string
	stat     counter.CodeStat
}

type fileList struct {
	data []*fileInfo
}

func getFiles(root string, filter string, files *fileList) error {
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if ok, _ := filepath.Match(filter, f.Name()); ok {
			fileinfo := &fileInfo{filename: path}
			files.data = append(files.data, fileinfo)
			//fmt.Println(fileinfo.filename)
		}
		return nil
	}

	return filepath.Walk(root, walkFunc)
}

func main() {
	filter := "*.go"
	root := "F:/dev/go_code/src/codecount"

	files := &fileList{}

	getFiles(root, filter, files)

	total_stat := &counter.CodeStat{}

	factory := counter.NewCodeCounterFactory()

	for _, v := range files.data {

		c, ok := factory.NewCounter("go")
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
