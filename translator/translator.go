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
	A        string
	e        []string
	comm     string
	floating bool
	tabs     string
}

var allRules []Rule

func main() {

	if len(os.Args) <= 1 {
		panic(errors.New("Not enough arguments! You need atleast the ABNF-file location. "))
	}

	abnf_path := os.Args[1]
	path_split := strings.Split(abnf_path, ".")
	filename := strings.Join(path_split[:len(path_split)-1], ".")

	// Open the reader
	f, err := os.Open(abnf_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
			allRules[len(allRules)-1].comm += "\n"
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
				line = strings.TrimSpace(comment)
			} else {
				line = comment
			}
			allRules = append(allRules, Rule{"", []string{}, line, false, ""})

			// Rule
		} else if isABNFRule(content) {
			content = ruleFmt(content, comment)

			line = content + "\t\t" + comment + "\n"

			// Additive line
		} else if strings.Contains(content, "=/") {
			i := len(allRules) - 1

			for len(allRules[i].A) <= 0 {
				i -= 1
			}

			// Add to all rules
			for ind, rul := range allRules {
				if utils.Strcmp(rul.A, strings.TrimSpace(allRules[i].A)) == 0 {
					e_ar := regexp.MustCompile(`[^"]/`).Split(contentFmt(strings.Split(content, "=/")[1]), -1)
					for _, e := range e_ar {
						found := false
						for _, e_p := range allRules[ind].e {
							if strings.Contains(e_p, e) {
								found = true
								break
							}
						}
						if !found {
							allRules[ind].e = append(allRules[ind].e, e)
						}
					}

					break
				}
			}
			continue

			// Floating content
		} else {
			var tabs string
			if isPEGRule(prev_line) {
				tabs = detTabs(strings.Index(prev_line, "<-") + 3)
			} else {
				tabs = detTabs(utils.CountPrefixSpace(prev_line))
			}
			content = contentFmt(content)
			allRules = append(allRules, Rule{allRules[len(allRules)-1].A, regexp.MustCompile(`[^"]/`).Split(content, -1), comment, true, tabs})
			content = tabs + content

			line = content + "\t\t" + comment + "\n"

		}

		prev_line = strings.Replace(line, "\t", "", -1)
	}

	// LEFT RECURSION REMOVAL
	var tempRules []Rule
	var groupRules map[string][]Rule
	groupRules = make(map[string][]Rule)

	for _, pRule := range allRules {
		groupRules[pRule.A] = append(groupRules[pRule.A], pRule)
	}

	for _, pRule := range allRules {
		group := groupRules[pRule.A]
		if len(group) > 1 && len(pRule.A) > 0 {
			A, AP := removeDirectLeftRecursion(group)

			if len(AP.A) > 0 && len(A.A) > 0 {
				if len(A.e) > 0 {
					tempRules = append(tempRules, A)
				}

				if len(AP.e) > 0 {
					tempRules = append(tempRules, AP)
				}
			} else {
				tempRules = append(tempRules, group...)
			}
			delete(groupRules, pRule.A)
		} else if len(group) != 0 {
			tempRules = append(tempRules, pRule)
		}
	}
	allRules = tempRules

	// INDIRECT LEFT RECURSION
	//allRules = removeIndirectLeftRecursion(allRules)

	// Write to output
	fo, err := os.Create(filename + ".peg")
	if err != nil {
		panic(err)
	}
	fo.WriteString(printRules(allRules))
	fo.Close()
	fmt.Println(filename+".peg", "generated")
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

	allRules = append(allRules, Rule{A, regexp.MustCompile(`[^"]/`).Split(e, -1), comment, false, ""})

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
	} else if !utils.IsUpper(word) && len(regexp.MustCompile("([a-z,A-Z]+)").FindString(word)) > 1 {
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
	return strings.Contains(content, `=`) && !strings.Contains(content, `"="`) && !strings.Contains(content, `=/`)
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
		word = strings.Replace(word, r.FindString(word), " "+rep+"? ", -1)
		//fmt.Println(word)

	}
	// one or more
	r, _ = regexp.Compile(`[1]\*[(](.+?)[)]|[1]\*(\S+)`)
	if len(r.FindString(word)) > 0 {
		rep := r.FindString(word)[2:] + "+"
		word = strings.Replace(word, r.FindString(word), rep, -1)
	}

	// zero or more
	r, err := regexp.Compile(`\*[(](.+?)[)]|\*([^\s\"]+)`)
	if err != nil {
		fmt.Println(err)
	}
	for len(r.FindString(word)) > 0 {
		rep := r.FindString(word)[1:]
		word = r.ReplaceAllString(word, rep+"*")
	}

	return word
}

func removeDirectLeftRecursion(prodRules []Rule) (Rule, Rule) {
	var A Rule
	var A_Prime Rule

	var lrRules []string
	var bRules []string

	for _, pRule := range prodRules {
		for _, comp := range pRule.e {
			if strings.HasPrefix(strings.TrimSpace(comp), pRule.A+" ") {
				lrRules = append(lrRules, strings.TrimSpace(comp))

				/*else if strings.HasPrefix(comp, regexp.MustCompile(`([^\"]+?)\?([\s]*)`+pRule.A+`([\s]*)`).FindString(comp)) {
					lrRules = append(lrRules, regexp.MustCompile(`([^\"]+?)\?([\s]*)`+pRule.A+`([\s]*)`).FindString(comp))

				} else if strings.HasPrefix(comp, regexp.MustCompile(`(.+?)\*([\s]*)`).FindString(comp)) {
					fmt.Println("HERE", regexp.MustCompile(`(.+?)\*`).FindString(comp))
					lrRules = append(lrRules, regexp.MustCompile(`\*([^\s\"]+)`+pRule.A+`([\s]+)`).FindString(comp))
				*/
			} else {
				bRules = append(bRules, strings.TrimSpace(comp))
			}
		}
	}

	if len(lrRules) > 0 {
		var ap string
		if len(bRules) > 0 {
			ap = prodRules[0].A + "'"
		} else {
			ap = prodRules[0].A
		}

		// A production
		A = Rule{prodRules[0].A, []string{}, prodRules[0].comm, prodRules[0].floating, prodRules[0].tabs}
		for _, beta := range bRules {
			A.e = append(A.e, beta+" "+ap)
		}

		// A' production
		A_Prime = Rule{ap, []string{"ε"}, "", false, ""}
		for _, alpha := range lrRules {
			A_Prime.e = append(A_Prime.e, strings.TrimPrefix(alpha, prodRules[0].A)+" "+A_Prime.A)
		}

		return A, A_Prime
	} else {
		return Rule{}, Rule{}
	}
}

func removeIndirectLeftRecursion(prodRules []Rule) []Rule {

	var tempRules []Rule

	for i := 0; i < len(prodRules); i++ {
		for j := 0; j < i-1; j++ {
			var temp []string
			for ind, e := range prodRules[i].e {
				if strings.HasPrefix(e, prodRules[j].A) {
					for _, s := range prodRules[j].e {
						temp = append(temp, s+strings.TrimPrefix(e, prodRules[j].A))
					}

					prodRules[i].e = append(append(prodRules[i].e[:ind], temp...), prodRules[i].e[ind+1:]...)
				}
			}
		}

		A, AP := removeDirectLeftRecursion([]Rule{prodRules[i]})

		if len(AP.A) > 0 && len(A.A) > 0 {
			tempRules = append(tempRules, A)
			tempRules = append(tempRules, AP)
		} else {
			tempRules = append(tempRules, prodRules[i])
		}
	}

	return tempRules
}

func printRules(all []Rule) string {
	var s []string

	for _, pRule := range all {
		for i, nt := range pRule.e {
			pRule.e[i] = strings.TrimSpace(nt)
		}

		if len(pRule.A) > 0 {
			if pRule.floating {
				if len(pRule.e) > 0 {
					s = append(s, pRule.tabs+strings.Join(pRule.e, " / "))

					if len(strings.TrimSpace(pRule.comm)) > 1 {
						s[len(s)-1] += "\t\t"
					}
					s[len(s)-1] += pRule.comm
				}
			} else {
				if len(pRule.e) > 0 {
					s = append(s, pRule.A+" <- "+strings.Join(pRule.e, " / "))

					if len(strings.TrimSpace(pRule.comm)) > 1 {
						s[len(s)-1] += "\t\t"
					}
					s[len(s)-1] += pRule.comm
				}
			}

		} else {
			s = append(s, pRule.comm)
		}
	}
	return strings.Join(s, "\n")
}
