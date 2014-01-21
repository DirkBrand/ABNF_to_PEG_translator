package utils

import (
	"errors"
	"strconv"
)

func ConvertNumberToHex(num string, typ rune) (int64, error) {
	var base int
	switch typ {
	case 'b':
		base = 2
		break
	case 'd':
		base = 10
		break
	case 'x':
		base = 16
		break
	default:
		return -1, errors.New(num + " cannot be formatted")
	}

	return strconv.ParseInt(num, base, 64)

}
