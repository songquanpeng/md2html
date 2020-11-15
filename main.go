package main

import (
	"io/ioutil"
	"log"
	"md2html/parser"
	"os"
)

func main() {
	for _, file := range os.Args[1:] {
		markdown, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Converting file %q.", file)
		root := parser.Parse(string(markdown))
		parser.PrintAST(root)

		log.Printf("Converted file saved at %q.", file)
	}
}
