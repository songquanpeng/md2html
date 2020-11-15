package converter

import (
	"fmt"
	"md2html/parser"
)

var htmlTemplate = `<html>
<head>
<link rel='stylesheet' href='https://cdn.jsdelivr.net/npm/mvp.css@1.6.2/mvp.min.css'>
<style>
.article {
    margin: auto;
max-width: 750px;
}
</style>
</head>
<body style="">%s</body>
</html>`

func Convert(markdown string, fullPage bool) (html string) {
	ast := parser.Parse(markdown)
	parser.PrintAST(ast)
	html = processArticleNode(ast)
	if fullPage {
		html = fmt.Sprintf(htmlTemplate, html)
	}
	return html
}

func processArticleNode(node *parser.Node) (html string) {
	for _, child := range node.Children {
		switch child.Type {
		case parser.TitleNode:
			html += processTitleNode(child)
		case parser.DividingLineNode:
			html += processDividingLineNode(child)
		case parser.ContentNode:
			content := processContentNode(child)
			html += fmt.Sprintf("<div>%s</div>\n", content)
		case parser.ListNode:
			html += processListNode(child)
		case parser.QuoteNode:
			html += processQuoteNode(child)
		case parser.CodeBlockNode:
			html += processCodeBlockNode(child)
		}
	}
	html = fmt.Sprintf("<div class='article'>\n%s\n</div>", html)
	return
}

func processTitleNode(node *parser.Node) (html string) {
	content := processContentNode(node.Children[0])
	level := int(node.Value[0])
	html = fmt.Sprintf("<h%d>%s</h%d>\n", level, content, level)
	return
}

func processDividingLineNode(node *parser.Node) (html string) {
	if node.Type == parser.DividingLineNode {
		html = "<hr>\n"
	}
	return
}

func processContentNode(node *parser.Node) (html string) {
	for _, child := range node.Children {
		switch child.Type {
		case parser.TextNode:
			html += fmt.Sprintf("%s", string(child.Value))
		case parser.ItalicNode:
			html += processRichTextNode(child, "i")
		case parser.BoldNode:
			html += processRichTextNode(child, "b")
		case parser.InlineCodeNode:
			html += processRichTextNode(child, "code")
		case parser.StrikethroughNode:
			html += processRichTextNode(child, "del")
		case parser.LinkNode:
			html += processLinkNode(child)
		case parser.ImageNode:
			html += processImageNode(child)
		case parser.ContentNode:
			html += processContentNode(child)
		}
	}
	html = fmt.Sprintf("%s", html)
	return
}

func processListNode(node *parser.Node) (html string) {

	return
}

func processQuoteNode(node *parser.Node) (html string) {
	content := processContentNode(node.Children[0])
	html = fmt.Sprintf("<q>%s</q>\n", content)
	return
}

func processCodeBlockNode(node *parser.Node) (html string) {
	content := string(node.Value)
	html = fmt.Sprintf("<pre><code>%s</code></pre>", content)
	return
}

func processRichTextNode(node *parser.Node, tag string) (html string) {
	content := processContentNode(node.Children[0])
	html = fmt.Sprintf("<%s>%s</%s>", tag, content, tag)
	return
}

func processLinkNode(node *parser.Node) (html string) {
	content := string(node.Children[0].Value)
	link := string(node.Value)
	html = fmt.Sprintf("<a href='%s'>%s</a>", link, content)
	return
}

func processImageNode(node *parser.Node) (html string) {
	content := string(node.Children[0].Value)
	link := string(node.Value)
	html = fmt.Sprintf("<img src='%s' alt='%s'/>", link, content)
	return
}
