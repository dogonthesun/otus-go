package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const escapeChar = '\\'

type prevChar struct {
	code   rune
	exists bool
}

func Unpack(str string) (string, error) {
	var sb strings.Builder
	var escapeMode bool

	p := prevChar{}

	for _, char := range str {
		if escapeMode {
			// escape only digits and escape character itself
			if !unicode.IsDigit(char) && char != escapeChar {
				return "", ErrInvalidString
			}
			p, escapeMode = prevChar{char, true}, false
			continue
		}

		if unicode.IsDigit(char) {
			// nothing to repeat
			if !p.exists {
				return "", ErrInvalidString
			}
			repeatCnt := int(char) - '0'
			for i := 0; i < repeatCnt; i++ {
				sb.WriteRune(p.code)
			}
			p = prevChar{}
			continue
		}

		if p.exists {
			sb.WriteRune(p.code)
		}

		if char == escapeChar {
			p, escapeMode = prevChar{}, true
			continue
		}

		p = prevChar{char, true}
	}

	if p.exists {
		sb.WriteRune(p.code)
	}

	// tailed escape
	if escapeMode {
		return "", ErrInvalidString
	}

	return sb.String(), nil
}
