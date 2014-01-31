package lexer

import (
	token "github.com/DirkBrand/ABNF_to_PEG_translator/Arithmetic_Parser/token"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	src []byte
	pos int
}

var ch rune

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src: src,
		pos: 0,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (this *Lexer) Scan() token.Token {
	if this.pos >= len(this.src) {
		return token.Token{token.EOF, token.EOF.String()}
	}
	size := 0

	ch, size = utf8.DecodeRune(this.src[this.pos:])
	this.pos += size

	if unicode.IsSpace(ch) {
		ch, size = utf8.DecodeRune(this.src[this.pos:])
		this.pos += size
	}

	if ch == '(' {
		return token.Token{token.LPAREN, "("}
	} else if ch == ')' {
		return token.Token{token.RPAREN, ")"}
	} else if ch == '*' {
		return token.Token{token.TIMES, "*"}
	} else if ch == '+' {
		return token.Token{token.PLUS, "+"}
	} else if unicode.IsDigit(ch) {
		buf := string(ch)
		for this.pos < len(this.src) {
			ch, size = utf8.DecodeRune(this.src[this.pos:])

			if unicode.IsDigit(ch) {
				buf += string(ch)
				this.pos += size
			} else {
				break
			}
		}
		return token.Token{token.NUM, buf}
	}

	return token.Token{token.EOF, token.EOF.String()}
}

func (this *Lexer) HasToken() bool {
	return this.pos < len(this.src)
}
