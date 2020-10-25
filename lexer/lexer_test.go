package lexer

import (
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

func TestTokenizeHeader(t *testing.T) {
	Tokenize(markdown1)
	n := 0
	for ; true; n++ {
		token := NextToken()
		//fmt.Printf("%d: <%s, %q>\n", n, TokenTypeName[token.Type], string(token.Value))
		if token.Type == EofToken {
			break
		}
	}
	expectedTokenNum := 19
	if n != expectedTokenNum {
		t.Errorf("There should be %d tokens", expectedTokenNum)
	}
}
