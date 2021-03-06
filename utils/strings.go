package utils

import (
	"strings"
	"unicode"
)

func PadZeros(a string, length int) string {
	s := a
	for len(s) != length {
		s = "0" + s
	}

	return s
}

func CountPrefixSpace(line string) (a int) {
	for _, c := range line {
		if unicode.IsSpace(c) {
			a += 1
		} else {
			break
		}
	}
	return a
}

func RemoveAllSpace(line string) (s string) {
	for _, c := range line {
		if !unicode.IsSpace(c) {
			s += string(c)
		}
	}
	return s
}

func ToCamelCase(array []string) string {
	var A string
	for _, nt := range array {
		A += strings.Title(nt)
	}
	return A
}

func IsUpper(s string) bool {
	flag := true
	for _, c := range s {
		if !unicode.IsUpper(c) {
			flag = false
		}
	}

	return flag
}

func Strcmp(a, b string) int {
	var min = len(b)
	if len(a) < len(b) {
		min = len(a)
	}
	var diff int
	for i := 0; i < min && diff == 0; i++ {
		diff = int(a[i]) - int(b[i])
	}
	if diff == 0 {
		diff = len(a) - len(b)
	}
	return diff
}
