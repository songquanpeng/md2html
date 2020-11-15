package parser

import (
	"fmt"
	"github.com/disiqueira/gotree"
	"log"
	"md2html/lexer"
)

type NodeType int8

const (
	NoneNode NodeType = iota
	ArticleNode
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
	"NoneNode",
	"ArticleNode",
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

func (node Node) String() (str string) {
	str += NodeTypeName[node.Type]
	switch node.Type {
	case TextNode:
		str += fmt.Sprintf(": %q", string(node.Value))
	case TitleNode:
		str += fmt.Sprintf(": %d", node.Value[0])
	case ImageNode:
		fallthrough
	case LinkNode:
		str += fmt.Sprintf(": %s", string(node.Value))
	case ListNode:
		str += ": "
		switch int(node.Value[0]) {
		case 0:
			str += "Unordered List"
		case 1:
			str += "Ordered List"
		case 2:
			str += "Uncompleted Task"
		case 3:
			str += "Completed Task"
		}
		str += fmt.Sprintf(" (level %d)", int(node.Value[1]))
	}
	return
}

// If you add any global variables, don't forget to progress them in func Parse(markdown string)!
var tokenBuffer []lexer.Token
var pos = 0
var tabCounter = 0

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
	return
}

func restoreToken() {
	if pos == 0 {
		log.Println("Warning: nothing to restore!")
	} else {
		pos--
	}
}

func nextTokenIs(tokenType lexer.TokenType) (yes bool) {
	token := getToken()
	yes = token.Type == tokenType
	restoreToken()
	return
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
	tokenBuffer = nil
	pos = 0
	tabCounter = 0
	lexer.Tokenize(markdown)
	root = parseArticle()
	preprocessAST(root)
	return
}

func preprocessAST(root *Node) {
	// Organize the list items as a tree, before that they are flatted.
	var newChildren []*Node
	for i := 0; i < len(root.Children); i++ {
		newChildren = append(newChildren, root.Children[i])
		if root.Children[i].Type == ListNode {
			start := i
			level := int(root.Children[i].Value[1])
			end := i
			for j := i + 1; j < len(root.Children); j++ {
				if root.Children[j].Type == ListNode {
					if int(root.Children[j].Value[1]) <= level {
						end = j
						break
					} else {
						end = j + 1
					}
				} else {
					end = j
					break
				}

			}
			if end > start+1 {
				target := root.Children[start:end]
				processListNode(target)
				i = end - 1
			}
		}
	}
	root.Children = newChildren
}

func processListNode(nodes []*Node) {
	// The first node is the root node.
	root := nodes[0]
	for i := 1; i < len(nodes); i++ {
		root.Children = append(root.Children, nodes[i])
		noMoreChild := true
		level := int(nodes[i].Value[1])
		for j := i + 1; j < len(nodes); j++ {
			if int(nodes[j].Value[1]) <= level {
				// Current one is a new child.
				noMoreChild = false
				// Firstly we should complete the previous child.
				processListNode(nodes[i:j])
				// Then we update i.
				i = j - 1
				break
			}
		}
		if noMoreChild {
			processListNode(nodes[i:])
			break
		}
	}
}

func parseArticle() (root *Node) {
	root = parseSectionList()
	root.Type = ArticleNode
	return
}

func parseSectionList() (root *Node) {
	node := Node{}
	root = &node
	for {
		token := getToken()
		restoreToken()
		current := &Node{}
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
			tabCounter = 0
			continue
		case lexer.TabToken:
			tabCounter++
			_ = getToken()
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
	// return -1 if not found, notice it doesn't not check the start point
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
		case lexer.LinkHeadToken:
			nodeType = LinkNode
		case lexer.ImageHeadToken:
			nodeType = ImageNode
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
		case lexer.LinkHeadToken:
			fallthrough
		case lexer.ImageHeadToken:
			if i+2 <= end &&
				(*tokens)[i+1].Type == lexer.TextToken &&
				(*tokens)[i+2].Type == lexer.LinkBodyToken {
				current.Type = getNodeTypeBySymToken((*tokens)[i])
				current.Value = (*tokens)[i+2].Value
				current.Children = append(current.Children, &Node{
					Type:     TextNode,
					Value:    (*tokens)[i+1].Value,
					Children: nil,
				})
				i += 3
			} else {
				// Not paired, fallback to text token and rerun this loop
				(*tokens)[i].Type = lexer.TextToken
				i--
				continue
			}
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
				// And rerun this loop.
				i--
				continue
			} else {
				// Paired
				current = constructRichTextNode(i, pos, tokens, getNodeTypeBySymToken((*tokens)[i]))
				// Don't forget to update i
				i = pos
			}
		default:
			// Fallback to text token.
			(*tokens)[i].Type = lexer.TextToken
			i--
			continue
		}
		root.Children = append(root.Children, current)
	}
	return
}

func constructRichTextNode(start, end int, tokens *[]lexer.Token, nodeType NodeType) (root *Node) {
	node := Node{}
	root = &node
	root.Type = nodeType
	subNode := constructContentNode(start+1, end-1, tokens)
	root.Children = append(root.Children, subNode)
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
	listType := 0
	switch token.Type {
	case lexer.UnorderedListToken:
		listType = 0
	case lexer.OrderedListToken:
		listType = 1
	case lexer.UncompletedTaskToken:
		listType = 2
	case lexer.CompletedTaskToken:
		listType = 3
	default:
		log.Println("Warning: unexpected token detected when processing list.")
	}
	listLevel := tabCounter
	tabCounter = 0
	root.Value = append(root.Value, rune(listType), rune(listLevel))
	// The first child of a list node is its content.
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
