package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goodplayer/vendordep/analyse"
)

type DepMain struct {
	Project Project
	Deps    []DepItem
}

type Project struct {
	GroupId        string
	Name           string
	ImportRootPath string
}

// groupId+name+version to find a dependency
type DepItem struct {
	GroupId        string
	Name           string
	Version        string
	ImportRootPath string
	VcsType        string
	VcsUrl         string
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "    get    -- get all deps: vendordep get")
		return
	}

	p, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Open(filepath.Join(p, "vendordep.json"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("read vendordep.json file error.", err)
	}

	var dep DepMain
	err = json.Unmarshal(data, &dep)
	if err != nil {
		log.Fatalln("unmarshal error.", err)
	}

	for _, v := range dep.Deps {
		log.Println("================================processing group:", v.GroupId, "name:", v.Name, "...")
		paths := strings.Split(v.ImportRootPath, "/")
		paths = append([]string{p, "vendor"}, paths...)
		path := filepath.Join(paths...)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			log.Println("mkdir:", path, "error.", err)
		}
		switch v.VcsType {
		case "git":
			processGit(v, path)
		default:
			log.Fatalln("unknown vcs type:", v.VcsType)
		}
		log.Println("================================processed group:", v.GroupId, "name:", v.Name)
	}

	// last: print unrecognized import
	allImports := analyse.GetImportPaths(p)
	existImports := make([]string, len(dep.Deps)+1)
	i := 0
	for _, v := range dep.Deps {
		existImports[i] = v.ImportRootPath
		i++
	}
	existImports[i] = dep.Project.ImportRootPath
	unimported := analyse.MergeUnimportedUrlPaths(allImports, existImports)

	if len(unimported) > 0 {
		fmt.Println("-------->>>> following imports are unrecognized. please check and add them to vendordep.json:")
		for _, v := range unimported {
			fmt.Println("\t", v)
		}
	} else {
		fmt.Println("-------->>>> process dependencies finished.")
	}
}

func processGit(v DepItem, dir string) {
	e, err := exec.LookPath("git")
	if err != nil {
		log.Fatalln("cannot lookup git executable file.")
	}
	cmd := exec.Command(e, "clone", v.VcsUrl, dir)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalln("git clone error.", err)
	}
	cmd2 := exec.Command(e, "checkout", v.Version)
	cmd2.Dir = dir
	cmd2.Env = os.Environ()
	cmd2.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	err = cmd2.Run()
	if err != nil {
		log.Fatalln("git checkout version error.", err)
	}
	err = os.RemoveAll(filepath.Join(dir, ".git"))
	if err != nil {
		log.Fatalln("remove .git folder error.", err)
	}
}
