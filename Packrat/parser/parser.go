package parser

import (
	"errors"
	"fmt"
	"strconv"
)

var rs []rune
var sym rune
var pos int

func Parse(t []rune) bool {
	rs = t
	pos = 0
	sym = rs[0]

	if !Mailbox() {
		return false
	}
	if sym == 127 {
		fmt.Println("Reached EOF")
	}

	return true
}

func accept(s rune) bool {
	if sym == s {
		getsym()
		return true
	}
	return false
}
func expect(s rune) (bool, error) {
	p := pos
	if accept(s) {
		return true, nil
	}
	err := errors.New("ERROR: Expected symbol " + string(s) + " but was " + string(sym))
	return back(p), err
}

func Mailbox() bool {
	fmt.Println("Mailbox")
	if NameAddr() {
		if sym != rune(127) {
			return false
		}
		return true
	} else if AddrSpec() {
		if sym != rune(127) {
			return false
		}
		return true
	}
	return false
}

func NameAddr() bool {
	fmt.Println("NameAddr")
	p := pos
	DisplayName()

	if !AngleAddr() {
		return back(p)
	}
	return true
}

func AngleAddr() bool {
	fmt.Println("AngleAddr")
	p := pos

	CFWS()

	if b, _ := expect('<'); !b {
		return false
	}

	if !AddrSpec() {
		return back(p)
	}
	if b, _ := expect('>'); !b {
		return false
	}

	CFWS()

	return true
}

func AddrSpec() bool {
	fmt.Println("AddrSpec")
	p := pos

	if !LocalPart() {
		return back(p)
	}
	if b, err := expect('@'); !b {
		fmt.Println(err)
		return false
	} else {
		fmt.Println("Accepted -> @")
	}
	if !Domain() {
		return back(p)
	}
	return true
}

func LocalPart() bool {
	fmt.Println("LocalPart")
	if DotAtom() {
		return true
	} else if QuotedString() {
		return true
	}
	return false
}

func Domain() bool {
	fmt.Println("Domain")
	p := pos
	if DotAtom() {
		return true
	} else if DomainLiteral() {
		return true
	}
	return back(p)
}

func DomainLiteral() bool {
	fmt.Println("DomainLiteral")
	CFWS()

	if b, err := expect('['); !b {
		fmt.Println(err)
		return false
	} else {
		fmt.Println("Accepted -> [")
	}

	for {
		p := pos
		FWS()
		if !Dtext() {
			back(p)
			break
		}
	}
	FWS()
	if b, err := expect(']'); !b {
		fmt.Println(err)
		return false
	} else {
		fmt.Println("Accepted -> ]")
	}
	CFWS()
	return true
}

func Dtext() bool {
	fmt.Println("Dtext")
	p := pos
	i := sym

	getsym()
	if i >= 33 && i <= 90 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	} else if i >= 94 && i <= 126 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	}
	return back(p)
}

func DisplayName() bool {
	fmt.Println("DisplayName")
	return Phrase()
}

func Word() bool {
	fmt.Println("Word")
	if Atom() {
		return true
	} else if QuotedString() {
		return true
	}
	return false
}

func Phrase() bool {
	fmt.Println("Phrase")
	if !Word() {
		return false
	}
	for Word() {
	}

	return true
}

func QuotedString() bool {
	fmt.Println("QuotedString")
	CFWS()
	if !DQUOTE() {
		return false
	}
	for {
		FWS()
		if !Qcontent() {
			break
		}
	}
	FWS()
	if !DQUOTE() {
		return false
	}
	CFWS()
	return true
}

func Atext() bool {
	fmt.Println("Atext")
	p := pos

	if ALPHA() {
		return true
	}
	if DIGIT() {
		return true
	}
	r := sym
	getsym()
	switch r {
	case '!':
		return true
	case '#':
		return true
	case '$':
		return true
	case '%':
		return true
	case '&':
		return true
	case '\'':
		return true
	case '*':
		return true
	case '+':
		return true
	case '-':
		return true
	case '/':
		return true
	case '=':
		return true
	case '?':
		return true
	case '^':
		return true
	case '_':
		return true
	case '`':
		return true
	case '{':
		return true
	case '|':
		return true
	case '}':
		return true
	case '~':
		return true
	default:
		return back(p)
	}
}

func ALPHA() bool {
	fmt.Println("ALPHA")
	p := pos
	i := sym

	getsym()
	if i >= 65 && i <= 90 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	} else if i >= 97 && i <= 122 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	} else if i >= 128 && i <= 165 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	}
	return back(p)
}

func Atom() bool {
	fmt.Println("Atom")
	CFWS()
	if !Atext() {
		return false
	}
	for Atext() {
	}
	CFWS()
	return true
}

func DotAtomText() bool {
	fmt.Println("DotAtomText")
	if !Atext() {
		return false
	}
	for Atext() {
	}
	p := pos
	for {
		if b, err := expect('.'); !b {
			fmt.Println(err)
			break
		} else {
			fmt.Println("Accepted -> .")
		}
		if !Atext() {
			back(p)
			break
		}
		for Atext() {
		}
	}

	return true
}

func DotAtom() bool {
	fmt.Println("DotAtom")
	CFWS()
	if !DotAtomText() {
		return false
	}
	CFWS()
	return true
}

func CFWS() bool {
	fmt.Println("CFWS")

	b := true
	FWS()
	if !Comment() {
		b = false
	}
	if b {
		p := pos
		for {
			FWS()
			if !Comment() {
				back(p)
				break
			}
		}

		FWS()
	}

	if b {
		return true
	}

	if FWS() {
		return true
	}

	return false
}

func FWS() bool {
	for {
		if !WSP() {
			break
		}
	}
	CRLF()

	if !WSP() {
		return false
	}
	for {
		if !WSP() {
			break
		}
	}
	return true
}

func Comment() bool {
	fmt.Println("Comment")

	if b, err := expect('('); !b {
		fmt.Println(err)
		return false
	}

	for {
		p := pos
		FWS()
		if !Ccontent() {
			back(p)
			break
		}
	}

	FWS()

	if b, err := expect(')'); !b {
		fmt.Println(err)
		return false
	}

	return true
}

func Ccontent() bool {
	fmt.Println("Ccontent")
	if Ctext() {
		return true
	} else if QuotedPair() {
		return true
	} else if Comment() {
		return true
	}
	return false
}

func Ctext() bool {
	fmt.Println("Ctext")
	p := pos

	i := sym

	getsym()
	if i >= 33 && i <= 39 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	} else if i >= 42 && i <= 91 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	} else if i >= 93 && i <= 122 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	}

	return back(p)
}

func QuotedPair() bool {
	fmt.Println("QuotedPair")

	if b, _ := expect('\\'); !b {
		return false
	}

	if VCHAR() {
		return true
	} else if WSP() {
		return true
	}
	return false
}

func Qtext() bool {
	fmt.Println("Qtext")
	p := pos

	i, err := strconv.ParseInt(fmt.Sprintf("%U", sym)[2:], 16, 32)
	if err != nil {
		panic(err)
		return false
	}
	getsym()
	if i == 33 {
		return true
	} else if i >= 35 && i <= 91 {
		return true
	} else if i >= 93 && i <= 122 {
		return true
	}

	return back(p)
}

func Qcontent() bool {
	fmt.Println("Qcontent")
	if Qtext() {
		return true
	} else if QuotedPair() {
		return true
	}
	return false
}

func CRLF() bool {
	p := pos

	i, err := strconv.ParseInt(fmt.Sprintf("%U", sym)[2:], 16, 32)
	if err != nil {
		panic(err)
		return false
	}
	getsym()
	if i == 13 {
		return true
	} else if i == 10 {
		return true
	}

	return back(p)
}

func WSP() bool {
	if SP() {
		return true
	} else if HTAB() {
		return true
	}
	return false
}

func SP() bool {
	p := pos

	i := sym

	getsym()
	if i == 32 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	}

	return back(p)
}

func HTAB() bool {
	p := pos

	i := sym

	getsym()
	if i == 9 {
		fmt.Println("Accepted ->", i, string(rs[pos-1]))
		return true
	}

	return back(p)
}

func DQUOTE() bool {
	fmt.Println("DQUOTE")
	p := pos

	r := sym
	getsym()

	if r == '"' {
		return true
	}

	return back(p)
}

func DIGIT() bool {
	fmt.Println("DIGIT")
	p := pos

	i, err := strconv.ParseInt(fmt.Sprintf("%U", sym)[2:], 16, 32)
	if err != nil {
		panic(err)
		return false
	}
	getsym()
	if i >= 48 && i <= 57 {
		fmt.Println("Accepted ->", string(rs[pos-1]))
		return true
	}

	return back(p)
}

func VCHAR() bool {
	p := pos

	i, err := strconv.ParseInt(fmt.Sprintf("%U", sym)[2:], 16, 32)
	if err != nil {
		panic(err)
		return false
	}
	getsym()
	if i >= 33 && i <= 126 {
		return true
	}

	return back(p)
}

func getsym() {
	if pos >= len(rs)-1 {
		fmt.Println(errors.New("ERROR: Reached EOF while reading " + string(sym)))
		return
	}
	pos += 1
	sym = rs[pos]
}

func back(p int) bool {
	fmt.Println("<BackTracking from", string(sym), "to", string(rs[p]), ">")
	pos = p
	sym = rs[pos]
	return false
}
