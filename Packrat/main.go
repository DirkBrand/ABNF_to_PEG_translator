package main

import (
	"errors"
	"fmt"
	"github.com/DirkBrand/ABNF_to_PEG_translator/Packrat/lexer"
	"github.com/DirkBrand/ABNF_to_PEG_translator/Packrat/parser"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(errors.New("Not enough arguments! You need atleast the file location. "))
	}

	filename := os.Args[1]

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	// Added delimiter
	src = append(src, 127)
	fmt.Println(src)

	b := parser.Parse(lexer.Scan(src))
	if b {
		fmt.Println("PROGRAM VALID")
	} else {
		fmt.Println("PARSE FAILED")

	}
}
