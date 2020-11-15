package lexer

import (
	"fmt"
	"testing"
)

const markdown1 = `
# h1
## h2
### h3
#### h4
##### h5
###### h6
`

func checkTokenNumber(t *testing.T, markdown string, expectedTokenNum int, showTokenStream bool) {
	Tokenize(markdown)
	n := 0
	for ; true; n++ {
		token := NextToken()
		if showTokenStream {
			fmt.Printf("%d: <%s, %q>\n", n, TokenTypeName[token.Type], string(token.Value))
		}
		if token.Type == EofToken {
			break
		}
	}
	if n != expectedTokenNum {
		t.Errorf("There should be %d tokens", expectedTokenNum)
	}
}

func TestTokenizeHeader(t *testing.T) {
	checkTokenNumber(t, markdown1, 19, false)
}

const markdown2 = `
---
hi
---`

func TestTokenizeDividingLine(t *testing.T) {
	checkTokenNumber(t, markdown2, 6, false)
}

const markdown3 = `
This is my [website](https://justsong.cn).
My avatar: ![alt text](https://justsong.cn/favicon.ico)
`

func TestTokenizeLinkAndImage(t *testing.T) {
	checkTokenNumber(t, markdown3, 12, false)
}

const markdown4 = `
* Item 1
* Item 2
	- [ ] Uncompleted task 1
	* Item 2b 
		- [x] Completed task 1
		* Item 2bb
			* Item 2bba
			* Item 2bbb

1. Item 1
1. Item 2 
1. Item 3
    1. Item 3a
	1. Item 3b
`

func TestTokenizeList(t *testing.T) {
	checkTokenNumber(t, markdown4, 55, false)
}
