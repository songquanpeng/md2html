package main

import (
	"io/ioutil"
	"log"
	"md2html/parser"
)

func main() {
	markdown, err := ioutil.ReadFile("./test/test.md")
	if err != nil {
		log.Fatal(err)
	}
	root := parser.Parse(string(markdown))
	parser.PrintAST(root)
}
