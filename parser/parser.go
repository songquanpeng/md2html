package parser

import (
	"fmt"
	"github.com/disiqueira/gotree"
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

var NodeTypeName = []string{
	"ArticleNode",
	"NoneNode",
	"ContentNode",
	"TextNode",
	"TitleNode",
	"DividingLineNode",
	"QuoteNode",
	"CodeBlockNode",
	"ListNode",
	"ItalicNode",
	"BoldNode",
	"InlineCodeNode",
	"StrikethroughNode",
	"LinkNode",
	"ImageNode",
}

type Node struct {
	Type     NodeType
	Value    []rune
	Children []*Node
}

func (node Node) String() string {
	return NodeTypeName[node.Type]
}

var tokenBuffer []lexer.Token
var pos = 0

func getToken() (token lexer.Token) {
	if pos == len(tokenBuffer) {
		// Noting in the buffer or all tokens are used.
		pos++
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
	//fmt.Printf("getToken() called with token ")
	//lexer.PrintToken(token)
	return
}

func restoreToken() {
	//fmt.Printf("restoreToken() called.\n")
	if pos == 0 {
		log.Println("Warning: nothing to restore!")
	} else {
		pos--
	}
}

func PrintAST(root *Node) {
	tree := gotree.New("Root")
	printASTHelper(root, &tree)
	fmt.Println(tree.Print())
}

func printASTHelper(astNode *Node, treeNode *gotree.Tree) {
	if astNode != nil {
		for _, child := range astNode.Children {
			subTree := (*treeNode).Add((*child).String())
			printASTHelper(child, &subTree)
		}
	}
}

func Parse(markdown string) (root *Node) {
	lexer.Tokenize(markdown)
	root = parseArticle()
	return
}

func parseArticle() (root *Node) {
	root = parseSectionList()
	return
}

func parseSectionList() (root *Node) {
	node := Node{}
	root = &node
	for {
		token := getToken()
		restoreToken()
		current := &Node{}
		//fmt.Printf("for loop start, current token: ")
		//lexer.PrintToken(token)
		switch token.Type {
		case lexer.TitleToken:
			current = parseTitle()
		case lexer.DividingLineToken:
			current = parseDividingLine()
		case lexer.CodeBlockToken:
			current = parseCodeBlock()
		case lexer.UncompletedTaskToken:
			fallthrough
		case lexer.CompletedTaskToken:
			fallthrough
		case lexer.UnorderedListToken:
			fallthrough
		case lexer.OrderedListToken:
			current = parseList()
		case lexer.QuoteToken:
			current = parseQuote()
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
			current = parseContent()
		}
		root.Children = append(root.Children, current)
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
	root.Children = append(root.Children, parseContent())
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
	// First we should retrieve all the tokens this content node need.
	var tokens []lexer.Token
	for token := getToken(); token.Type != lexer.NewlineToken && token.Type != lexer.EofToken; token = getToken() {
		tokens = append(tokens, token)
	}
	return constructContentNode(0, len(tokens)-1, &tokens)
}

func constructContentNode(start, end int, tokens *[]lexer.Token) (root *Node) {
	node := Node{}
	root = &node
	root.Type = ContentNode
	if tokens == nil {
		log.Println("Warning: content node is blank!")
		return
	}
	// return -1 if not found, do not check the start point
	findThisTypeToken := func(t lexer.TokenType, start int) (pos int) {
		for pos = start + 1; pos < len(*tokens); pos++ {
			if (*tokens)[pos].Type == t {
				return pos
			}
		}
		pos = -1
		return
	}

	getNodeTypeBySymToken := func(token lexer.Token) (nodeType NodeType) {
		switch token.Type {
		case lexer.SingleStarToken:
			fallthrough
		case lexer.SingleUnderscoreToken:
			nodeType = ItalicNode
		case lexer.DoubleStarToken:
			fallthrough
		case lexer.DoubleUnderscoreToken:
			nodeType = BoldNode
		case lexer.SingleBacktickToken:
			nodeType = InlineCodeNode
		case lexer.DoubleTildeToken:
			nodeType = StrikethroughNode
		default:
			nodeType = TextNode
		}
		return
	}
	for i := start; i <= end; i++ {
		current := &Node{}
		switch (*tokens)[i].Type {
		case lexer.TextToken:
			node := Node{}
			current = &node
			node.Type = TextNode
			node.Value = (*tokens)[i].Value
		case lexer.DoubleStarToken:
			fallthrough
		case lexer.DoubleUnderscoreToken:
			fallthrough
		case lexer.SingleStarToken:
			fallthrough
		case lexer.SingleUnderscoreToken:
			fallthrough
		case lexer.SingleBacktickToken:
			fallthrough
		case lexer.DoubleTildeToken:
			pos := findThisTypeToken((*tokens)[i].Type, i)
			if pos == -1 {
				// Not paired, change its type.
				(*tokens)[i].Type = lexer.TextToken
				// And rerun this cycle.
				i--
				continue
			} else {
				// Paired
				current = constructRichTextNode(i, pos, tokens, getNodeTypeBySymToken((*tokens)[i]))
				// Don't forget to update i
				i = pos
				continue
			}
		}
		root.Children = append(root.Children, current)
	}
	return
}

func constructRichTextNode(start, end int, tokens *[]lexer.Token, nodeType NodeType) (root *Node) {
	node := Node{}
	root = &node
	root.Type = nodeType
	root.Children[0] = constructContentNode(start+1, end-1, tokens)
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
	root.Children = append(root.Children, parseContent())
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
	root.Children = append(root.Children, parseContent())
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
