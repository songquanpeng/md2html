package parser

import (
	"io/ioutil"
	"log"
	"md2html/lexer"
	"testing"
)

func TestGetAndRestoreToken(t *testing.T) {
	markdown, err := ioutil.ReadFile("./test/test.md")
	if err != nil {
		log.Fatal(err)
	}
	lexer.Tokenize(string(markdown))
	for {
		token := getToken()
		lexer.PrintToken(token)
		if token.Type == lexer.EofToken {
			break
		}
	}
}
