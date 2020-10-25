package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"md2html/lexer"
)

func main() {
	markdown, err := ioutil.ReadFile("./test/basic.md")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print("Reading markdown...\n")
	//fmt.Print(string(markdown), "\n")
	//fmt.Print("Processing...\n")
	lexer.Tokenize(string(markdown))
	for i := 0; true; i++ {
		token := lexer.NextToken()
		fmt.Printf("%d: <%s, %q>\n", i, lexer.TokenTypeName[token.Type], string(token.Value))
		if token.Type == lexer.EofToken {
			break
		}
	}
}
