package lexer

import (
	"unicode"
)

type TokenType int8

const (
	EofToken TokenType = iota
	TextToken
	NewlineToken
	TabToken
	SingleStarToken
	DoubleStarToken
	SingleUnderscoreToken
	DoubleUnderscoreToken
	SingleBacktickToken
	CodeBlockToken
	DoubleTildeToken
	TitleToken
	UnorderedListToken
	OrderedListToken
	QuoteToken
	DividingLineToken
	UncompletedTaskToken
	CompletedTaskToken
	LinkHeadToken
	ImageHeadToken
	LinkBodyToken
)

var TokenTypeName = []string{
	"EofToken",
	"TextToken",
	"NewlineToken",
	"TabToken",
	"SingleStarToken",
	"DoubleStarToken",
	"SingleUnderscoreToken",
	"DoubleUnderscoreToken",
	"SingleBacktickToken",
	"CodeBlockToken",
	"DoubleTildeToken",
	"TitleToken",
	"UnorderedListToken",
	"OrderedListToken",
	"QuoteToken",
	"DividingLineToken",
	"UncompletedTaskToken",
	"CompletedTaskToken",
	"LinkHeadToken",
	"ImageHeadToken",
	"LinkBodyToken",
}

type Token struct {
	Type  TokenType
	Value []rune
}

var input []rune
var pos = 0
var lastTokenType = NewlineToken
var tokenQueue []Token

func Tokenize(markdown string) {
	input = []rune(markdown)
}

func nextIsSameTo(c rune) bool {
	if pos+1 >= len(input) {
		return false
	}
	return c == input[pos+1]
}

func isSpaceBehind() bool {
	if pos+1 >= len(input) {
		return true
	}
	return input[pos+1] == ' '
}

func isNumDotSpace() bool {
	if unicode.IsDigit(input[pos]) {
		return len(input) > pos+2 && input[pos+1] == '.' && input[pos+2] == ' '
	}
	return false
}

func isTaskSymbol() (yes, completed bool) {
	yes = false
	if len(input) > pos+2 && input[pos] == '[' {
		completed = input[pos+1] != ' '
		if input[pos+2] == ']' {
			yes = true
		}
	}
	return
}

func countSymbol(c rune) (n int) {
	n = 0
	for pos+n < len(input) && input[pos+n] == c {
		n++
	}
	return
}

func getCodeBlockStartEnd() (start, end int) {
	start = pos
	end = pos
	// Skip the language name.
	for ; start < len(input) && input[start] != '\n'; start++ {
	}
	start++
	for end = start; end+2 < len(input); end++ {
		if input[end] == '`' && input[end+1] == '`' && input[end+2] == '`' {
			break
		}
	}
	return
}

func NextToken() (token Token) {
	if len(tokenQueue) != 0 {
		token = tokenQueue[0]
		tokenQueue = tokenQueue[1:]
	} else {
		textToken, otherToken := nextToken()
		if len(textToken.Value) != 0 {
			token = textToken
			tokenQueue = append(tokenQueue, otherToken)
		} else {
			token = otherToken
		}
		lastTokenType = otherToken.Type
	}
	return
}

func nextToken() (textToken, otherToken Token) {
	textToken.Type = TextToken
	for {
		if pos >= len(input) {
			otherToken.Type = EofToken
			return
		}
		c := input[pos]
		if len(textToken.Value) == 0 && lastTokenType == NewlineToken {
			switch c {
			case '#':
				n := countSymbol(c)
				otherToken.Type = TitleToken
				otherToken.Value = append(otherToken.Value, rune(n))
				pos += n
				if input[pos] == ' ' {
					pos++
				}
				return
			case '\t':
				otherToken.Type = TabToken
				pos++
				return
			case '\n':
				otherToken.Type = NewlineToken
				pos++
				return
			case '-':
				fallthrough
			case '*':
				if isSpaceBehind() {
					otherToken.Type = UnorderedListToken
					pos += 2
					yes, completed := isTaskSymbol()
					if yes {
						pos += 2
						if isSpaceBehind() {
							pos += 2
							if completed {
								otherToken.Type = CompletedTaskToken
							} else {
								otherToken.Type = UncompletedTaskToken
							}
							return
						}
						pos -= 2
					}
					return
				} else { // Consider if this is a dividing line
					if nextIsSameTo(c) {
						pos++
						if nextIsSameTo(c) {
							pos++
							otherToken.Type = DividingLineToken
							return
						}
						pos--
					}
				}
			case '>':
				if isSpaceBehind() {
					otherToken.Type = QuoteToken
					pos += 2
					return
				}
			case '`':
				if nextIsSameTo(c) {
					pos++
					if nextIsSameTo(c) {
						pos += 2
						otherToken.Type = CodeBlockToken
						start, end := getCodeBlockStartEnd()
						otherToken.Value = input[start:end]
						pos = end + 3
						return
					}
					pos--
				}
			case '\r':
				fallthrough
			case ' ':
				n := countSymbol(c)
				if n >= 2 {
					pos += n
					otherToken.Type = TabToken
					return
				} else {
					pos++
				}
			}
			if isNumDotSpace() {
				otherToken.Type = OrderedListToken
				return
			}
		}
		// Update c because pos maybe updated due to black symbol.
		c = input[pos]

		// Now we have to return the text token before the below token.
		switch c {
		case '*':
			if nextIsSameTo(c) {
				pos += 2
				otherToken.Type = DoubleStarToken
				otherToken.Value = []rune("**")
			} else {
				pos += 1
				otherToken.Type = SingleStarToken
				otherToken.Value = []rune("*")
			}
			return
		case '_':
			if nextIsSameTo(c) {
				pos += 2
				otherToken.Type = DoubleUnderscoreToken
				otherToken.Value = []rune("__")
			} else {
				pos += 1
				otherToken.Type = SingleUnderscoreToken
				otherToken.Value = []rune("_")
			}
			return
		case '~':
			if nextIsSameTo(c) {
				pos += 2
				otherToken.Type = DoubleTildeToken
				otherToken.Value = []rune("~~")
				return
			}
		case '`':
			otherToken.Type = SingleBacktickToken
			otherToken.Value = []rune("`")
			pos++
			return
		case '!':
			if nextIsSameTo('[') {
				pos += 2
				otherToken.Type = ImageHeadToken
				return
			}
		case '[':
			pos++
			otherToken.Type = LinkHeadToken
			return
		case ']':
			if nextIsSameTo('(') {
				pos += 2
				for i := pos; i < len(input) && input[i] != '\n'; i++ {
					if input[i] == ')' {
						otherToken.Type = LinkBodyToken
						otherToken.Value = input[pos:i]
						pos = i + 1
						return
					}
				}
				pos -= 2
			}
		case '\n':
			otherToken.Type = NewlineToken
			pos++
			return
		}
		pos++
		if c != '\r' {
			textToken.Value = append(textToken.Value, c)
		}
	}
}
