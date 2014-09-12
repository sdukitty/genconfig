package main

import (
	"flag"
	// "github.com/deepglint/glog"
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type DyncValue struct {
	Name string
}

func main() {

	var root = flag.String("path", ".", "root path_required!")
	flag.Parse()

	// walk through root
	absroot, _ := filepath.Abs(*root)
	err := filepath.Walk(absroot, func(curpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking file %s: %s", *root, err)
			return err
		}

		if nil == info {
			log.Printf("Error info: %s", err)
			return nil
		}

		curdir := path.Dir(curpath)
		if curdir != absroot {
			return nil
		}

		name := path.Base(curpath)
		if strings.HasPrefix(name, ".") {
			return nil
		}

		if info.IsDir() {
			// generate and write config file
			GenConfigFile(curpath)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking root: %s", err)
	}
}

func GenConfigFile(path string) {
	name := filepath.Base(path)
	var dyncvalue DyncValue = DyncValue{Name: name}
	tmpl, err := template.ParseFiles("./config.txt")
	if err != nil {
		log.Printf("Error reading template: %s", err)
		return
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, dyncvalue)
	if err != nil {
		log.Printf("Error parsing template: %s", err)
		return
	}
	tarfile := filepath.Join(path, "config.txt")
	log.Println(tarfile)
	if PathExist(tarfile) {
		os.Remove(tarfile)
	}
	file, err := os.OpenFile(tarfile, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Printf("Error writing to target file: %s", err)
		return
	}
	file.WriteString(doc.String())
	defer file.Close()
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
