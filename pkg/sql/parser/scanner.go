package parser

import "fmt"

/*
NUMBER         → DIGIT+ ( "." DIGIT+ )? ;
STRING         → "\"" <any char except "\"">* "\"" ;
IDENTIFIER     → ALPHA ( ALPHA | DIGIT )* ;
ALPHA          → "a" ... "z" | "A" ... "Z" | "_" ;
DIGIT          → "0" ... "9" ;
KEYWORDS       → "SELECT" | "FROM" | "WHERE" | "GROUP BY" | "OFFSET" | "LIMIT"
*/

func Scan(s string) ([]*Token, error) {
	tokens := []*Token{}
	i := 0

	for !isAtEnd(s, i+1) {
		// first check if at reserved keyword or symbol
		for j := i + 6; j >= i+1; j-- {
			if len(s) < j {
				continue
			}
			word := s[i:j]
			if ttype, ok := keywordLookup[word]; ok {
				tokens = append(tokens, simpleToken(ttype, word))
				i += len(word)
				continue
			}
		}

		switch c := s[i : i+1]; c {
		case " ", "\n", "\t":
			i++
		case "'":
			token, err := scanStr(s, i+1)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i += len(token.Lexeme) + 2
		default:
			fmt.Println(s[i : i+6])
			return nil, fmt.Errorf("unknown character '%s'", c)
		}
	}
	return tokens, nil
}

func isAtEnd(s string, i int) bool {
	return len(s) <= i
}

func scanStr(s string, start int) (*Token, error) {
	i := start
	for !isAtEnd(s, i) {
		switch c := s[i : i+1]; c {
		case "'":
			return newToken(STRING, s[start:i], s[start:i]), nil
		default:
			i++

		}
	}
	return nil, fmt.Errorf("reached end of input parsing string")
}
