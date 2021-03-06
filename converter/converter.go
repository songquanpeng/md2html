package converter

import (
	"fmt"
	"md2html/parser"
	"os"
	"strings"
)

func Convert(markdown string, fullPage bool) (html string) {
	ast := parser.Parse(markdown)
	if os.Getenv("MODE") == "debug" {
		parser.PrintAST(ast)
	}
	html = processArticleNode(ast)
	if fullPage {
		html = fmt.Sprintf(HtmlTemplate, Style, html)
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
	if len(node.Children) == 0 {
		return
	}
	if int(node.Children[0].Value[0])%2 == 0 {
		// unordered list
		html = "<ul>%s</ul>"
	} else {
		html = "<ol>%s</ol>"
	}
	content := ""
	for _, child := range node.Children {
		content += processSubListNode(child)
	}
	html = fmt.Sprintf(html, content)
	return
}

func processSubListNode(node *parser.Node) (html string) {
	content := processContentNode(node.Children[0])
	i := strings.Index(content, ". ")
	if i >= 0 && i < 5 {
		i += 2
		content = content[i:]
	}
	inputTag := ""
	if int(node.Value[0]) == 2 {
		inputTag = "<input disabled type='checkbox'>"
	} else if int(node.Value[0]) == 3 {
		inputTag = "<input checked disabled type='checkbox'>"
	}
	subListContent := ""
	if len(node.Children) > 1 {
		subList := ""
		if int(node.Children[1].Value[0])%2 == 0 {
			subList = "<ul>%s</ul>"
		} else {
			subList = "<ol>%s</ol>"
		}
		for _, child := range node.Children[1:] {
			subListContent += fmt.Sprintf(subList, processSubListNode(child))
		}
	}
	html += fmt.Sprintf("<li>%s%s%s</li>", inputTag, content, subListContent)
	return
}

func processQuoteNode(node *parser.Node) (html string) {
	content := ""
	for _, child := range node.Children {
		content += processContentNode(child)
	}
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
