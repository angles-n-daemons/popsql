package scanner

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// scanner.go implements the logic for tokenizing a sql query.
// For the lexical grammar, consult grammar.bnf.

var whitespace = []byte{' ', '\n', '\t'}

var Debug = false

func Scan(s string) ([]*Token, error) {
	tokens := []*Token{}
	i := 0

outside:
	for !isAtEnd(s, i) {
		// first check if at reserved keyword or symbol
		minLen := int(math.Min(float64(i+6), float64(len(s))))
		upper := strings.ToUpper(s[i:minLen])
		for j := minLen - i; j >= 1; j-- {
			word := upper[:j]
			if ttype, ok := keywordLookup[word]; ok {
				tokens = append(tokens, simpleToken(ttype, word))
				i += len(word)
				continue outside
			}
		}

		switch {
		case slices.Contains(whitespace, s[i]):
			i++
		case s[i] == '"':
			token, err := scanStr(s, i+1)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i += len(token.Lexeme) + 2
		case isNumeric(s[i]):
			token, err := scanNum(s, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i += len(token.Lexeme)
		case isLetter(s[i]):
			token, err := scanIdentifier(s, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i += len(token.Lexeme)
		default:
			return nil, fmt.Errorf("unknown character '%c'", s[i])
		}
	}
	if Debug {
		debugLexemes := []string{}
		debugTypes := []string{}
		for _, token := range tokens {
			llen := len(token.Lexeme)
			tlen := len(token.Type.String())
			if llen <= tlen {
				debugLexemes = append(debugLexemes, token.Lexeme+strings.Repeat(" ", tlen-llen))
				debugTypes = append(debugTypes, token.Type.String())
			} else {

				debugLexemes = append(debugLexemes, token.Lexeme)
				debugTypes = append(debugTypes, token.Type.String()+strings.Repeat(" ", llen-tlen))
			}
		}
		fmt.Println("Tokens: ")
		fmt.Println("[ " + strings.Join(debugLexemes, ", ") + " ]")
		fmt.Println("[ " + strings.Join(debugTypes, ", ") + " ]")
	}
	return tokens, nil
}

func isNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

func isLetter(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isAtEnd(s string, i int) bool {
	return len(s) <= i
}

func scanStr(s string, start int) (*Token, error) {
	i := start + 1
	for !isAtEnd(s, i) {
		switch s[i] {
		case '"':
			return newToken(STRING, s[start:i], s[start:i]), nil
		default:
			i++

		}
	}
	return nil, fmt.Errorf("reached end of input scanning string")
}

func scanNum(s string, start int) (*Token, error) {
	i := start
	dotFound := false
loop:
	for !isAtEnd(s, i) {
		switch {
		case s[i] == '.':
			if dotFound {
				return nil, fmt.Errorf("found second decimal in numeric value")
			}
			dotFound = true
			i++
		case isNumeric(s[i]):
			i++
		default:
			break loop
		}
	}
	num, err := strconv.ParseFloat(s[start:i], 64)
	if err != nil {
		return nil, err
	}
	return newToken(NUMBER, s[start:i], num), nil
}

func scanIdentifier(s string, start int) (*Token, error) {
	i := start + 1
	for !isAtEnd(s, i) && (isLetter(s[i]) || isNumeric(s[i]) || s[i] == '_') {
		i++
	}
	return newToken(IDENTIFIER, s[start:i], s[start:i]), nil
}

func PrintTokens(tokens []*Token) {
	for _, token := range tokens {
		fmt.Printf("%s: %s\n", token.Type, token.Lexeme)
	}
}
