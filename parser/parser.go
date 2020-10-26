package parser

import (
	"log"
	"md2html/lexer"
)

type NodeType int8

const (
	ArticleNode NodeType = iota
	NoneNode
	ContentNode
	TextNode
	TitleNode
	DividingLineNode
	QuoteNode
	CodeBlockNode
	ListNode
	ItalicNode
	BoldNode
	InlineCodeNode
	StrikethroughNode
	LinkNode
	ImageNode
)

type Node struct {
	Type     NodeType
	Value    []rune
	children [3]*Node
	Next     *Node
}

var tokenBuffer []lexer.Token
var pos = 0

func getToken() (token lexer.Token) {
	if pos == len(tokenBuffer) {
		// Noting in the buffer or all tokens are used.
		token = lexer.NextToken()
		tokenBuffer = append(tokenBuffer, token)
	} else {
		token = tokenBuffer[pos]
		pos++
	}
	if pos > 5 {
		// Remove the most outdated token
		tokenBuffer = tokenBuffer[1:]
		pos--
	}
	return
}

func restoreToken() {
	if pos == 0 {
		log.Println("Warning: nothing to restore!")
	} else {
		pos--
	}
}

func Parse() (root *Node) {
	root = parseArticle()
	return
}

func parseArticle() (root *Node) {
	root = parseSectionList()
	return
}

func parseSectionList() (root *Node) {
	current := root
	token := getToken()
	restoreToken()
	for {
		switch token.Type {
		case lexer.TitleToken:
			current.Next = parseTitle()
		case lexer.DividingLineToken:
			current.Next = parseDividingLine()
		case lexer.CodeBlockToken:
			current.Next = parseCodeBlock()
		case lexer.UncompletedTaskToken:
			fallthrough
		case lexer.CompletedTaskToken:
			fallthrough
		case lexer.UnorderedListToken:
			fallthrough
		case lexer.OrderedListToken:
			current.Next = parseList()
		case lexer.QuoteToken:
			current.Next = parseQuote()
		case lexer.NewlineToken:
			_ = getToken()
			token = getToken()
			restoreToken()
			continue
		case lexer.TabToken:
			// TODO: process the tab token.
			_ = getToken()
			token = getToken()
			restoreToken()
			continue
		case lexer.EofToken:
			return
		default:
			current.Next = parseContent()
		}
		if current.Next != nil {
			current = current.Next
		}
	}
}

func parseTitle() (root *Node) {
	token := getToken()
	if token.Type != lexer.TitleToken {
		log.Println("Error: not a title token!")
	}
	node := Node{}
	root = &node
	root.Type = TitleNode
	root.Value = token.Value
	root.Next = parseContent()
	return
}

func parseDividingLine() (root *Node) {
	token := getToken()
	if token.Type != lexer.DividingLineToken {
		log.Println("Error: not a dividing line token!")
	}
	node := Node{}
	root = &node
	root.Type = DividingLineNode
	return
}

func parseContent() (root *Node) {
	// Markdown is not a CFG language.
	node := Node{}
	root = &node
	root.Type = ContentNode
	// First we should retrieve all the tokens this content node need.
	var tokens []lexer.Token
	for token := getToken(); token.Type != lexer.NewlineToken && token.Type != lexer.EofToken; token = getToken() {
		tokens = append(tokens, token)
	}
	current := root
	l := len(tokens)
	if tokens == nil {
		log.Println("Warning: content node is blank!")
		return
	}

	// return -1 if not found, do not check the start point
	findThisTypeToken := func(t lexer.TokenType, start int) (pos int) {
		for pos = start + 1; pos < len(tokens); pos++ {
			if tokens[pos].Type == t {
				return pos
			}
		}
		pos = -1
		return
	}

	for i := 0; i < l; i++ {
		switch tokens[i].Type {
		case lexer.TextToken:
			node := Node{}
			current = &node
			node.Type = TextNode
			node.Value = tokens[i].Value
			current = current.Next
		case lexer.SingleStarToken:
			pos := findThisTypeToken(tokens[i].Type, i)
			if pos == -1 {
				// Not paired, change its type.
				tokens[i].Type = lexer.TextToken
				// And rerun this cycle.
				i--
				continue
			} else {
				// TODO: checkpoint
			}
		case lexer.SingleUnderscoreToken:

		}

	}

	root.children[0] = parseText()
	root.children[1] = parseRichText()
	root.children[2] = parseText()
	return
}

// If the next token is not a valid text token, this function will return an empty one.
func parseText() (root *Node) {
	node := Node{}
	root = &node
	root.Type = TextNode
	token := getToken()
	if token.Type != lexer.TextToken {
		restoreToken()
	} else {
		root.Value = token.Value
	}
	return
}

// If the next token is not a valid rich text token, this function will return a nil.
func parseRichText() (root *Node) {
	token := getToken()
	restoreToken()
	switch token.Type {
	case lexer.SingleStarToken:
		fallthrough
	case lexer.SingleUnderscoreToken:
		parseItalic()
	case lexer.DoubleStarToken:
		fallthrough
	case lexer.DoubleUnderscoreToken:
		parseBold()
	case lexer.DoubleTildeToken:
		parseStrikeThrough()
	case lexer.LinkHeadToken:
		parseLink()
	case lexer.ImageHeadToken:
		parseImage()
	case lexer.SingleBacktickToken:
		parseInlineCode()
	}
}

func parseItalic() (root *Node) {
	node := Node{}
	root = &node
	return
}

func parseBold() (root *Node) {
	node := Node{}
	root = &node
	return
}

func parseInlineCode() (root *Node) {
	node := Node{}
	root = &node
	return
}

func parseStrikeThrough() (root *Node) {
	node := Node{}
	root = &node
	return
}

func parseLink() (root *Node) {
	node := Node{}
	root = &node
	return
}

func parseImage() (root *Node) {
	token := getToken()
	if token.Type != lexer.ImageHeadToken {
		log.Println("Error: not a image head token!")
	}
	node := Node{}
	root = &node
	// TODO
	return
}

func parseQuote() (root *Node) {
	token := getToken()
	if token.Type != lexer.QuoteToken {
		log.Println("Error: not a quote token!")
	}
	node := Node{}
	root = &node
	root.Type = QuoteNode
	root.Next = parseContent()
	return
}

func parseList() (root *Node) {
	token := getToken()
	if token.Type != lexer.UnorderedListToken && token.Type != lexer.OrderedListToken &&
		token.Type != lexer.UncompletedTaskToken && token.Type != lexer.CompletedTaskToken {
		log.Println("Error: not a list token!")
	}
	node := Node{}
	root = &node
	root.Type = ListNode
	root.Value = []rune(lexer.TokenTypeName[token.Type])
	root.Next = parseContent()
	return
}

func parseCodeBlock() (root *Node) {
	token := getToken()
	if token.Type != lexer.CodeBlockToken {
		log.Println("Error: not a code block token!")
	}
	node := Node{}
	root = &node
	root.Type = CodeBlockNode
	root.Value = token.Value
	return
}
