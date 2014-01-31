package lexer

import (
	"unicode/utf8"
)

var ch rune

func Scan(src []byte) (rs []rune) {
	pos := 0
	for pos < len(src) {
		ch, size := utf8.DecodeRune(src[pos:])
		pos += size
		rs = append(rs, ch)
	}

	return rs
}
