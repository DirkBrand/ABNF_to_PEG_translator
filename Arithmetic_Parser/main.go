package main

import (
	"errors"
	"fmt"
	"github.com/DirkBrand/ABNF_to_PEG_translator/Arithmetic_Parser/lexer"
	"github.com/DirkBrand/ABNF_to_PEG_translator/Arithmetic_Parser/parser"
	"github.com/DirkBrand/ABNF_to_PEG_translator/Arithmetic_Parser/token"
	"os"
)

var symbs []token.Token

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(errors.New("Not enough arguments! You need atleast the file location. "))
	}

	filename := os.Args[1]

	l, err := lexer.NewLexerFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	for l.HasToken() {
		symbs = append(symbs, l.Scan())
		//fmt.Println(symbs[len(symbs)-1])
	}

	parser.Parse(symbs)
	fmt.Println("PROGRAM VALID")

}
