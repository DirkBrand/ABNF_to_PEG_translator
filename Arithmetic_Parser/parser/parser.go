package parser

import (
	"errors"
	"fmt"
	"github.com/DirkBrand/ABNF_to_PEG_translator/Arithmetic_Parser/token"
	"os"
)

var tokens []token.Token
var sym token.Token
var pos int

func Parse(t []token.Token) {
	tokens = t
	pos = 0
	sym = tokens[0]

	Additive()
}

func accept(s token.Type) bool {
	if sym.Type == s {
		getsym()
		return true
	}
	return false
}
func expect(s token.Type) {
	if accept(s) {
		return
	}
	fmt.Println(errors.New("ERROR: Expected symbol " + s.String() + " but was " + sym.Type.String()))
	os.Exit(1)
}

func Additive() {
	Multiplicative()

	if sym.Type == token.PLUS {
		expect(token.PLUS)
		Additive()
	}
}

func Multiplicative() {
	Primary()
	if sym.Type == token.TIMES {
		expect(token.TIMES)
		Multiplicative()
	}

}

func Primary() {
	if sym.Type == token.LPAREN {
		expect(token.LPAREN)
		Additive()
		expect(token.RPAREN)
	} else {
		Decimal()
	}
}

func Decimal() {
	expect(token.NUM)
}

func getsym() {
	if pos == len(tokens)-1 {
		fmt.Println(errors.New("ERROR: Reached EOF while reading " + sym.Type.String()))
		os.Exit(1)

	}
	pos += 1
	sym = tokens[pos]
}
