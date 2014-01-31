package token

import (
	"fmt"
	"io"
)

type Type int

const (
	PLUS Type = iota
	TIMES
	LPAREN
	RPAREN
	NUM
	EOF
)

func (t Type) String() string {
	s := ""
	switch t {
	case PLUS:
		s = "+"
		break

	case TIMES:
		s = "*"
		break

	case LPAREN:
		s = "("
		break

	case RPAREN:
		s = ")"
		break

	case NUM:
		s = "NUM"
		break

	case EOF:
		s = fmt.Sprint(io.EOF)
	}

	return s
}

type Token struct {
	Type
	Val string
}

func (t Token) String() string {
	return fmt.Sprintf("%v : %v", t.Type.String(), t.Val)
}
