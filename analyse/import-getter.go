package analyse

import (
	"container/list"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetImportPaths(name string) []string {
	l := list.New()
	goThroughDir(name, l)
	arr := make([][]*ast.ImportSpec, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		arr[i] = e.Value.([]*ast.ImportSpec)
		i++
	}
	m := make(map[string]string)
	for _, v := range arr {
		for _, vv := range v {
			m[vv.Path.Value] = vv.Path.Value
		}
	}
	strList := make([]string, len(m))
	i = 0
	for k, _ := range m {
		strList[i] = k[1 : len(k)-1]
		i++
	}
	sort.Strings(strList)
	return strList
}

func goThroughDir(name string, l *list.List) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatalln("go through dir error.", err)
	}
	fi, err := os.Stat(name)
	if err != nil {
		log.Fatalln("go through dir error when stat file.", err)
	}
	if fi.IsDir() {
		//log.Println("processing dir:", name)
		names, err := f.Readdirnames(0)
		if err != nil {
			log.Fatalln("read dir names error.", err)
		}
		for _, v := range names {
			newPath := filepath.Join(name, v)
			goThroughDir(newPath, l)
		}
	} else {
		if strings.HasSuffix(name, ".go") {
			processFile(name, l)
		}
	}
}

func processFile(name string, l *list.List) {
	f, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
	if err != nil {
		log.Fatalln("process file error.", err)
	}
	l.PushBack(f.Imports)
}
