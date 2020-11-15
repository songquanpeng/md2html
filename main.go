package main

import (
	"io/ioutil"
	"log"
	"md2html/converter"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	for _, file := range os.Args[1:] {
		markdown, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Converting file %q.", file)
		html := converter.Convert(string(markdown), true)
		convertedFilename := strings.TrimSuffix(file, filepath.Ext(file))
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
}
