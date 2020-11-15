package main

import (
	"io/ioutil"
	"log"
	"md2html/converter"
	"os"
	"path/filepath"
	"strings"
)

func ConvertFile(path string) {
	markdown, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Converting file %q.", path)
	html := converter.Convert(string(markdown), true)
	convertedFilename := strings.TrimSuffix(path, filepath.Ext(path))
	convertedFilename += ".html"
	convertedFile, err := os.Create(convertedFilename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = convertedFile.WriteString(html)
	if err != nil {
		log.Fatal(err)
	}
	_ = convertedFile.Close()
	log.Printf("Converted file saved at %q.", convertedFilename)
}

func main() {
	var files []string
	for _, path := range os.Args[1:] {
		fi, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if matched, err := filepath.Match("*.md", filepath.Base(path)); err != nil {
					return err
				} else if matched {
					files = append(files, path)
				}
				return nil
			})
		case mode.IsRegular():
			files = append(files, path)
		}
	}
	for _, file := range files {
		ConvertFile(file)
	}
}
