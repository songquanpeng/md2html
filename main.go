package main

import "md2html/lexer"

func main()  {
	lexer.Tokenize("##hi\n**Bold**")
}