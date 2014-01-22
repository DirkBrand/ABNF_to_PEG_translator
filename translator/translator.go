/*

Copyright (c) 2013, Dirk Brand
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted
provided that the following conditions are met:

 * Redistributions of source code must retain the above copyright notice, this list of
   conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of
   conditions and the following disclaimer in the documentation and/or other materials provided
   with the distribution.

THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND ANY EXPRESS OR IMPLIED
WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS
BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT
OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	utils "github.com/DirkBrand/ABNF_to_PEG_translator/utils"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	A    string
	e    []string
	comm string
}

var allRules []Rule

func main() {

	if len(os.Args) <= 1 {
		panic(errors.New("Not enough arguments! You need atleast the ABNF-file location. "))
	}

	abnf_path := os.Args[1]
	filename := strings.Split(abnf_path, ".")[0]

	var s string

	// Open the reader
	f, err := os.Open(abnf_path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	prev_line := ""

	for {
		line, err := r.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			s += "\n"
			continue
		}
		line = commentFmt(line)

		var content string
		var comment string
		i := strings.LastIndex(line, "#")
		if i >= 0 {
			content = line[:i]
			comment = strings.TrimSpace(line[i:])
		} else {
			content = line
		}

		// Only Comment
		if len(content) == 0 {
			var tabs string
			tabs = detTabs(strings.LastIndex(prev_line, "#"))
			comment = tabs + "\t\t" + comment
			if strings.LastIndex(prev_line, "#") <= 0 {
				line = strings.TrimSpace(comment) + "\n"
			} else {
				line = comment + "\n"
			}

			// Rule
		} else if isABNFRule(content) {
			content = ruleFmt(content, comment)

			line = content + "\t\t" + comment + "\n"

			// Floating content
		} else {
			var tabs string
			if isPEGRule(prev_line) {
				tabs = detTabs(strings.Index(prev_line, "<-") + 3)
			} else {
				tabs = detTabs(utils.CountPrefixSpace(prev_line))
			}
			content = contentFmt(content)
			allRules[len(allRules)-1].e = append(allRules[len(allRules)-1].e, regexp.MustCompile(`[^"]/`).Split(content, -1)...)
			content = tabs + content

			line = content + "\t\t" + comment + "\n"

		}

		prev_line = strings.Replace(line, "\t", "", -1)

		s += line
	}

	var tempRules []Rule
	var someFix bool
	for _, pRule := range allRules {
		//fmt.Println(pRule)
		A, AP := removeDirectLeftRecursion(pRule)
		tempRules = append(tempRules, A)
		if len(AP.A) > 0 {
			tempRules = append(tempRules, AP)
			someFix = true
		}

	}
	allRules = tempRules

	if someFix {
		// Write to output
		fo, err := os.Create(filename + ".RR.peg")
		if err != nil {
			panic(err)
		}
		fo.WriteString(printRules(allRules))
		fo.Close()
		fmt.Println("Left-Recursion free PEG generated")
	}

	// Write to output
	fo, err := os.Create(filename + ".peg")
	if err != nil {
		panic(err)
	}
	fo.WriteString(s)
	fo.Close()
	fmt.Println("PEG generated")
}

func contentFmt(content string) string {
	e := ""

	content = strings.TrimSpace(content)
	content = regexFmt(content)

	if strings.Contains(content, "/") {
		r, _ := regexp.Compile(`[^"]/`)
		e_arr := r.Split(content, -1)
		for _, sent := range e_arr {
			if len(strings.TrimSpace(sent)) > 0 {
				for _, word := range strings.Split(sent, " ") {
					if len(word) > 0 {
						e += wordFmt(strings.TrimSpace(word)) + " "
					}
				}
				e = strings.TrimSpace(e) + " / "
			}
		}
		if !strings.HasSuffix(content, "/") {
			e = strings.TrimSuffix(e, " / ")
		}
	} else {
		for _, word := range strings.Split(content, " ") {
			if len(word) > 0 {
				e += wordFmt(strings.TrimSpace(word)) + " "
			}
		}
		if strings.HasSuffix(content, "/") {
			e += "/"
		}
	}
	return strings.TrimSpace(e)
}

func ruleFmt(line string, comment string) string {
	var s string
	arr := strings.Split(line, "=")
	A_arr := strings.Split(strings.TrimSpace(arr[0]), "-")
	A := utils.ToCamelCase(A_arr)
	e := contentFmt(arr[1])

	allRules = append(allRules, Rule{A, regexp.MustCompile(`[^"]/`).Split(e, -1), comment})

	s += A + " <- " + e
	return s
}

func wordFmt(word string) string {
	var s string
	// Numeric terminal
	if strings.HasPrefix(word, `%`) {
		typ := rune(word[1])
		numbers := word[2:]
		// Range
		if strings.Contains(numbers, `-`) {
			num1 := strings.Split(numbers, "-")[0]
			num2 := strings.Split(numbers, "-")[1]

			a, err1 := utils.ConvertNumberToHex(num1, typ)
			b, err2 := utils.ConvertNumberToHex(num2, typ)
			if err1 != nil || err2 != nil {
				fmt.Println(err1, " ---- ", err2)
			}
			s += "'\\u" + utils.PadZeros(strconv.FormatInt(a, 16), 4) + "'-'\\u" + utils.PadZeros(strconv.FormatInt(b, 16), 4) + "'"

		} else {
			typ := rune(word[1])
			num1 := word[2:]
			a, err1 := utils.ConvertNumberToHex(num1, typ)
			if err1 != nil {
				fmt.Println(err1)
			}
			s += "'\\u" + utils.PadZeros(strconv.FormatInt(a, 16), 4) + "'"

		}

		// String Literal
	} else if match, _ := regexp.MatchString(`"(.+)"`, word); match {
		s = word

		// Non-Terminal
	} else if !utils.IsUpper(word) && len(word) > 1 {
		s = utils.ToCamelCase(strings.Split(word, "-"))
	} else {
		s = word
	}

	return s
}

func commentFmt(line string) string {
	return strings.Replace(line, ";", "#", -1)
}

func isABNFRule(content string) bool {
	return strings.Contains(content, `=`) && !strings.Contains(content, `"="`)
}

func isPEGRule(content string) bool {
	return strings.Contains(content, `<-`) && !strings.Contains(content, `"<-"`)
}

func detTabs(count int) (s string) {
	for i := 0; i < count; i++ {
		s += " "
	}
	return s
}

func regexFmt(word string) string {
	// Optional regex
	r, _ := regexp.Compile(`\[([^\"]+?)\]`)
	for len(r.FindString(word)) > 0 {
		//fmt.Println(r.FindString(word))
		rep := r.FindString(word)[1 : len(r.FindString(word))-1]
		if strings.Contains(rep, " ") {
			rep = "(" + rep + ")"
		}
		word = strings.Replace(word, r.FindString(word), rep+"?", -1)
		//fmt.Println(word)

	}
	// one or more
	r, _ = regexp.Compile(`[1]\*[(](.+?)[)]|[1]\*(\S+)`)
	if len(r.FindString(word)) > 0 {
		//fmt.Println(r.FindString(word), " ------> ", word)
		rep := r.FindString(word)[2:] + "+"
		word = strings.Replace(word, r.FindString(word), rep, -1)
		//fmt.Println(word + "\n")
	}

	// zero or more
	r, err := regexp.Compile(`\*[(](.+?)[)]|\*([^\s\"]+)`)
	if err != nil {
		fmt.Println(err)
	}
	for len(r.FindString(word)) > 0 {
		//fmt.Println(r.FindString(word))
		rep := r.FindString(word)[1:]
		word = r.ReplaceAllString(word, rep+"*")
		//fmt.Println(word)
	}

	return word
}

func removeDirectLeftRecursion(prodRule Rule) (Rule, Rule) {
	var A Rule
	var A_Prime Rule

	var lrRules []string
	var bRules []string

	for _, comp := range prodRule.e {
		if strings.HasPrefix(strings.TrimSpace(comp), prodRule.A+" ") {
			lrRules = append(lrRules, strings.TrimSpace(comp))
		} else {
			bRules = append(bRules, strings.TrimSpace(comp))
		}
	}

	if len(lrRules) > 0 {
		// A production
		A = Rule{prodRule.A, []string{}, prodRule.comm}
		for _, beta := range bRules {
			A.e = append(A.e, beta+" "+prodRule.A+"'")
		}

		// A' production
		A_Prime = Rule{prodRule.A + "'", []string{"ε"}, ""}
		for _, alpha := range lrRules {
			A_Prime.e = append(A_Prime.e, strings.TrimPrefix(alpha, prodRule.A)+" "+A_Prime.A)
		}

		return A, A_Prime
	} else {
		return prodRule, Rule{}
	}
}

func printRules(all []Rule) string {
	var s []string
	for _, pRule := range all {
		s = append(s, pRule.A+" <- "+strings.Join(pRule.e, " /")+"\t\t"+pRule.comm+"\n")
	}
	return strings.Join(s, "\n")
}
