package lexer

type TokenType int8

const (
	EOF_TOKEN TokenType = iota
	TEXT_TOKEN
	NEXT_LINE_TOKEN
	ITALIC_STAR_LEFT_TOKEN
	ITALIC_STAR_RIGHT_TOKEN
	ITALIC_UNDERLINE_LEFT_TOKEN
	ITALIC_UNDERLINE_RIGHT_TOKEN
	BOLD_STAR_LEFT_TOKEN
	BOLD_STAR_RIGHT_TOKEN
	BOLD_UNDERLINE_STAR_LEFT_TOKEN
	BOLD_UNDERLINE_STAR_RIGHT_TOKEN
	DELETE_LEFT_TOKEN
	DELETE_RIGHT_TOKEN
	TITLE_1_TOKEN
	TITLE_2_TOKEN
	TITLE_3_TOKEN
	TITLE_4_TOKEN
	TITLE_5_TOKEN
	TITLE_6_TOKEN
	DIVIDING_LINE_TOKEN
	INLINE_CODE_LEFT_TOKEN
	INLINE_CODE_RIGHT_TOKEN
	BLOCK_CODE_LEFT_TOKEN
	BLOCK_CODE_RIGHT_TOKEN
	TASK_UNCOMPLETED_TOKEN
	TASK_COMPLETED_TOKEN
	LINK_LEFT_TOKEN
	LINK_MIDDLE_TOKEN
	LINK_RIGHT_TOKEN
	IMAGE_LEFT_TOKEN
	IMAGE_MIDDLE_TOKEN
	IMAGE_RIGHT_TOKEN
)

type Token struct {
	Type TokenType
	Value []rune
}

var input []rune
var pos = 0

func Tokenize(markdown string) {
	input = []rune(markdown)
}

func NextToken() (token Token) {
	for {
		if pos >= len(input){
			token.Type = EOF_TOKEN
			return
		}
		nextLine, generateToken := processNextLine()
		if generateToken {
			token.Type = NEXT_LINE_TOKEN
			return
		}
		if nextLine {
			if level := processTitle(); level != 0 {
				if level > 0 {
					nextLine = false
				}
				switch level {
				case 1:
					token.Type = TITLE_1_TOKEN
				case 2:
					token.Type = TITLE_2_TOKEN
				case 3:
					token.Type = TITLE_3_TOKEN
				case 4:
					token.Type = TITLE_4_TOKEN
				case 5:
					token.Type = TITLE_5_TOKEN
				default:
					token.Type = TITLE_6_TOKEN
				}
				return
			}
			if processDividingLine() {
				token.Type = DIVIDING_LINE_TOKEN
				return
			}
		}
	}
}

func processNextLine() (nextLine, generateToken bool){
	nextLine = false
	generateToken = false
	if input[pos] == '\n' {
		nextLine = true
		pos++
		for pos < len(input) && input[pos] == '\n' {
			generateToken = true
			pos++
		}
	}
	return
}

func processTitle()  (level int8) {
	level = 0
	for pos < len(input) && input[pos] == '#' {
		level++
	}
	return
}

func processText()  {

}

func processDividingLine() (yes bool){
	yes = false
	oldPos := pos
	for pos < len(input) && input[pos] == '-' {
		pos++
	}
	if pos < len(input) && input[pos] == '\n' {
		yes = true
	} else {
		pos = oldPos
	}
	return
}