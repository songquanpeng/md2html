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
	TripleBacktickToken
	DoubleTildeToken
	Title1Token
	Title2Token
	Title3Token
	Title4Token
	Title5Token
	Title6Token
	UnorderedListToken
	OrderedListToken
	QuoteToken
	DividingLineToken
	UncompletedTaskToken
	CompletedTaskToken
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
	"TripleBacktickToken",
	"DoubleTildeToken",
	"Title1Token",
	"Title2Token",
	"Title3Token",
	"Title4Token",
	"Title5Token",
	"Title6Token",
	"UnorderedListToken",
	"OrderedListToken",
	"QuoteToken",
	"DividingLineToken",
	"UncompletedTaskToken",
	"CompletedTaskToken",
}

type Token struct {
	Type TokenType
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
	if pos + 1 >= len(input) {
		return false
	}
	return c == input[pos + 1]
}

func isSpaceBehind()  bool {
	if pos + 1 >= len(input) {
		return true
	}
	return input[pos + 1] == ' '
}

func isNumDotSpace() bool {
	if unicode.IsDigit(input[pos]) {
		return len(input) > pos + 2 && input[pos + 1] == '.' && input[pos + 2] == ' '
	}
	return false
}

func isTaskSymbol() (yes, completed bool) {
	yes = false
	if len(input) > pos + 2 && input[pos] == '[' {
		completed = input[pos + 1] != ' '
		if input[pos + 2] == ']' {
			yes = true
		}
	}
	return
}

func countSharp()  (n int) {
	n = 0
	for pos + n < len(input) && input[pos + n] == '#' {
		n++
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
		if pos >= len(input){
			otherToken.Type = EofToken
			return
		}
		c := input[pos]
		if len(textToken.Value) == 0 && lastTokenType == NewlineToken {
			switch c {
			case '#':
				n := countSharp()
				switch n {
				case 2:
					otherToken.Type = Title2Token
				case 3:
					otherToken.Type = Title3Token
				case 4:
					otherToken.Type = Title4Token
				case 5:
					otherToken.Type = Title5Token
				case 6:
					otherToken.Type = Title6Token
				default:
					otherToken.Type = Title1Token
				}
				pos += n
				if input[pos] == ' ' {
					pos ++
				}
				return
			case '\t':
				otherToken.Type = TabToken
				pos ++
				return
			case '\n':
				otherToken.Type = NewlineToken
				pos ++
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
						pos ++
						if nextIsSameTo(c) {
							pos ++
							otherToken.Type = DividingLineToken
							return
						}
						pos --
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
					pos ++
					if nextIsSameTo(c) {
						pos ++
						otherToken.Type = TripleBacktickToken
						return
					}
					pos --
				}
			case '\r':
				fallthrough
			case ' ':
				pos ++
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
			} else {
				pos += 1
				otherToken.Type = SingleStarToken
			}
			return
		case '_':
			if nextIsSameTo(c) {
				pos += 2
				otherToken.Type = DoubleUnderscoreToken
			} else {
				pos += 1
				otherToken.Type = SingleUnderscoreToken
			}
			return
		case '~':
			if nextIsSameTo(c) {
				pos += 2
				otherToken.Type = DoubleTildeToken
				return
			}
		case '`':
			otherToken.Type = SingleBacktickToken
			pos ++
			return
		case '\n':
			otherToken.Type = NewlineToken
			pos ++
			return
		}
		pos ++
		if c != '\r' {
			textToken.Value = append(textToken.Value, c)
		}
	}
}
