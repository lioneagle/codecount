package main

import (
	"counter"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileInfo struct {
	fileName  string
	shortName string
	ext       string
	stat      counter.CodeStat
}

type FileList []*FileInfo

func (f FileList) GetFileNameMaxLen() (fullNameMaxLen, shortNameMaxLen int) {
	for _, v := range f {
		if len(v.fileName) > fullNameMaxLen {
			fullNameMaxLen = len(v.fileName)
		}

		if len(v.shortName) > shortNameMaxLen {
			shortNameMaxLen = len(v.shortName)
		}
	}
	return fullNameMaxLen, shortNameMaxLen
}

func GetFiles(root string, filters []string, files *FileList) error {
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		for _, v := range filters {
			if ok, _ := filepath.Match(v, f.Name()); ok {
				fileinfo := &FileInfo{fileName: path, shortName: filepath.Base(path)}
				fileinfo.ext = strings.ToLower(filepath.Ext(path)[1:])
				*files = append(*files, fileinfo)
			}
		}
		return nil
	}

	return filepath.Walk(root, walkFunc)
}

type ExtMapToCodeType struct {
	maps map[string]string
}

func NewExtMapToCodeType() *ExtMapToCodeType {
	return &ExtMapToCodeType{maps: make(map[string]string)}
}

func (c *ExtMapToCodeType) BindFiltersToCodeType(filters string, codetype string) {
	ext := strings.Split(filters, ";")

	for _, v := range ext {
		if v == "" {
			continue
		}
		c.maps[strings.ToLower(filepath.Ext(v)[1:])] = strings.ToLower(codetype)
	}
}

type RunConfig struct {
	root          string
	filter        string
	showEachFile  bool
	showShortName bool
	sortStat      bool
	sortField     string
	sortReverse   bool
	csvOutput     bool
	csvFileName   string

	exts []*string
}

func (runConfig *RunConfig) Parse(codeConfigs []CodeConfig) {
	flag.StringVar(&runConfig.root, "path", ".", "path for code")
	flag.StringVar(&runConfig.filter, "filter", "*.cpp;*.cxx;*.hpp;*.hxx;*.c++;*.cc;*.c;*.h;*.go;*.java;*.erl;*.hrl", "file filters")
	flag.BoolVar(&runConfig.showEachFile, "show", false, "show each file stat")
	flag.BoolVar(&runConfig.showShortName, "short", true, "show file name without path")
	flag.BoolVar(&runConfig.sortStat, "sort", true, "sort stat result")
	flag.StringVar(&runConfig.sortField, "sortfield", "code", "set sort field: fullname, shortname, total, code, comment, blank, comment-percent")
	flag.BoolVar(&runConfig.sortReverse, "reverse", true, "sort reverse")
	flag.BoolVar(&runConfig.csvOutput, "csv", true, "enable to output csv file")
	flag.StringVar(&runConfig.csvFileName, "csvfile", "result.csv", "csv file name")

	flag.Parse()

	runConfig.exts = make([]*string, 0)
	for _, v := range codeConfigs {
		runConfig.exts = append(runConfig.exts, flag.String(v.codeType, v.filters, v.filtersDesc))
	}
}

func (runConfig *RunConfig) Check() bool {
	_, err := os.Stat(runConfig.root)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		fmt.Printf("ERROR: path \"%s\" is not exist", runConfig.root)
		return false
	}
	fmt.Printf("ERROR: path \"%s\" is invalid", runConfig.root)
	return false
}

type CodeConfig struct {
	codeType    string
	filters     string
	filtersDesc string
}

type CodeTypeStat struct {
	filenum int
	stat    counter.CodeStat
}

type AllStats struct {
	totalStat     counter.CodeStat
	codeTypeStats map[string]*CodeTypeStat
	codeTypeOrder []string
	maxPrefixLen  int
	totalFiles    int
}

func NewAllStats() *AllStats {
	return &AllStats{codeTypeStats: make(map[string]*CodeTypeStat)}
}

func (allStats *AllStats) getMaxPrefixLen(files FileList, runConfig *RunConfig) {
	allStats.maxPrefixLen = 0
	if runConfig.showEachFile {
		fullNameMaxLen, shortNameMaxLen := files.GetFileNameMaxLen()
		if runConfig.showShortName {
			allStats.maxPrefixLen = shortNameMaxLen
		} else {
			allStats.maxPrefixLen = fullNameMaxLen
		}
	}

	for k, v := range allStats.codeTypeStats {
		if v.filenum > 0 {
			str := fmt.Sprintf("total %d %s files", v.filenum, k)
			if len(str) > allStats.maxPrefixLen {
				allStats.maxPrefixLen = len(str)
			}
		}
	}

	str := fmt.Sprintf("total %d files", len(files))
	if len(str) > allStats.maxPrefixLen {
		allStats.maxPrefixLen = len(str)
	}
}

func (allStats *AllStats) AddCodeType(codetype string) {
	allStats.codeTypeStats[codetype] = &CodeTypeStat{}
	allStats.codeTypeOrder = append(allStats.codeTypeOrder, codetype)
}

func (allStats *AllStats) AddStat(codetype string, stat *counter.CodeStat) {

	allStats.totalStat.Add(stat)

	codetypeStat, ok := allStats.codeTypeStats[codetype]
	if !ok {
		codetypeStat = &CodeTypeStat{}
		allStats.codeTypeStats[codetype] = codetypeStat
	}

	codetypeStat.filenum++
	codetypeStat.stat.Add(stat)

	allStats.totalFiles++
}

func (allStats *AllStats) Print() (ret string) {

	for _, v := range allStats.codeTypeOrder {
		codeTypeStat, _ := allStats.codeTypeStats[v]
		if codeTypeStat.filenum > 0 {
			str := fmt.Sprintf("total %d %s files", codeTypeStat.filenum, v)
			ret += fmt.Sprintf("%s:  ", str)
			ret += PrintIdent(allStats.maxPrefixLen - len(str))
			ret += fmt.Sprintf("%s\n", codeTypeStat.stat.String())
		}
	}

	str := fmt.Sprintf("total %d files", allStats.totalFiles)
	ret += fmt.Sprintf("%s:  ", str)
	ret += PrintIdent(allStats.maxPrefixLen - len(str))
	ret += fmt.Sprintf("%s\n", allStats.totalStat.String())
	return ret
}

func (allStats *AllStats) WriteToCsvFile(w *csv.Writer) {
	for _, v := range allStats.codeTypeOrder {
		codeTypeStat, _ := allStats.codeTypeStats[v]
		if codeTypeStat.filenum > 0 {
			line := []string{fmt.Sprintf("total %d %s files", codeTypeStat.filenum, v)}
			line = append(line, codeTypeStat.stat.StringSlice()...)
			w.Write(line)
		}
	}

	line := []string{fmt.Sprintf("total %d files", allStats.totalFiles)}
	line = append(line, allStats.totalStat.StringSlice()...)
	w.Write(line)
}

func main() {

	codeConfigs := []CodeConfig{
		{"cpp", "*.cpp;*.cxx;*.hpp;*.hxx;*.c++;*.cc", "extions for c/c++ files"},
		{"c", "*.c;*.h", "extions for c files"},
		{"go", "*.go", "extions for go files"},
		{"java", "*.java", "extions for java files"},
		{"erlang", "*.erl;*.hrl", "extions for erlang files"},
	}

	runConfig := RunConfig{}
	runConfig.Parse(codeConfigs)
	if !runConfig.Check() {
		return
	}

	allStats := NewAllStats()
	extMapToCodeType := NewExtMapToCodeType()
	for _, v := range codeConfigs {
		extMapToCodeType.BindFiltersToCodeType(v.filters, v.codeType)
		allStats.AddCodeType(v.codeType)
	}

	files := FileList{}
	GetFiles(runConfig.root, strings.Split(runConfig.filter, ";"), &files)

	Run(files, extMapToCodeType, allStats)
	OutputResult(files, &runConfig, allStats)
}

func Run(files FileList, extMapToCodeType *ExtMapToCodeType, allStats *AllStats) {
	factory := counter.NewCodeCounterFactory()

	for _, v := range files {
		codeType, ok := extMapToCodeType.maps[v.ext]
		if !ok {
			log.Printf("ERROR: unknown code type for %s", v.fileName)
			continue
		}

		c, ok := factory.NewCounter(codeType)
		if !ok {
			log.Printf("ERROR: cannot get codecounter for %s", v.fileName)
			continue
		}

		stat, ok := counter.ParseFile(c, v.fileName)
		if !ok {
			log.Printf("ERROR: parse file %s failed", v.fileName)
		}

		v.stat = stat
		allStats.AddStat(codeType, &stat)

	}
}

func PrintResult(files FileList, runConfig *RunConfig, allStats *AllStats) (ret string) {
	allStats.getMaxPrefixLen(files, runConfig)
	SortResult(files, runConfig)
	if runConfig.showEachFile {
		for _, v := range files {
			if runConfig.showShortName {
				ret += fmt.Sprintf("%s:  ", v.shortName)
				ret += PrintIdent(allStats.maxPrefixLen - len(v.shortName))
			} else {
				ret += fmt.Sprintf("%s:  ", v.fileName)
				ret += PrintIdent(allStats.maxPrefixLen - len(v.fileName))
			}
			ret += fmt.Sprintf("%s\n", v.stat.String())
		}

		ret += "\n"
	}

	ret += allStats.Print()
	ret += "\n"

	return ret
}

func SortResult(files FileList, runConfig *RunConfig) {
	type compareFunc func(i, j int) bool

	sortobject := map[string]compareFunc{
		"fullname":        func(i, j int) bool { return files[i].fileName < files[j].fileName },
		"shortname":       func(i, j int) bool { return files[i].shortName < files[j].shortName },
		"total":           func(i, j int) bool { return files[i].stat.Total < files[j].stat.Total },
		"code":            func(i, j int) bool { return files[i].stat.Code < files[j].stat.Code },
		"comment":         func(i, j int) bool { return files[i].stat.Comment < files[j].stat.Comment },
		"blank":           func(i, j int) bool { return files[i].stat.Blank < files[j].stat.Blank },
		"comment-percent": func(i, j int) bool { return files[i].stat.CommentPercent() < files[j].stat.CommentPercent() },
	}

	if runConfig.sortStat {
		name := strings.ToLower(runConfig.sortField)
		if name == "name" {
			if runConfig.showShortName {
				name = "shortname"
			} else {
				name = "fullname"
			}
		}
		v, ok := sortobject[name]
		if ok {
			if runConfig.sortReverse {
				sort.Slice(files, func(i, j int) bool { return !v(i, j) })
			} else {
				sort.Slice(files, v)
			}
		}
	}
}

func OutputResult(files FileList, runConfig *RunConfig, allStats *AllStats) {
	ret := PrintResult(files, runConfig, allStats)
	fmt.Printf("%s", ret)

	if runConfig.csvOutput {
		OutputToCsvFile(files, runConfig, allStats)
	}
}

func OutputToCsvFile(files FileList, runConfig *RunConfig, allStats *AllStats) {
	file, err := os.OpenFile(runConfig.csvFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Printf("ERROR: cannot open csv file %s to write", runConfig.csvFileName)
		return
	}
	//file.WriteString("test")
	defer file.Close()

	w := csv.NewWriter(file)
	w.Write([]string{"FileName", "Total", "Code", "Comment", "Blank", "CommentPercent"})
	for _, v := range files {
		line := []string{}
		if runConfig.showShortName {
			line = append(line, v.shortName)
		} else {
			line = append(line, v.fileName)
		}
		line = append(line, v.stat.StringSlice()...)
		w.Write(line)
	}

	allStats.WriteToCsvFile(w)

	w.Flush()
}

func PrintIdent(num int) (ret string) {
	for i := 0; i < num; i++ {
		ret += " "
	}
	return ret
}
