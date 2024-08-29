package scanner

import "fmt"

/*
NUMBER         → DIGIT+ ( "." DIGIT+ )? ;
STRING         → "\"" <any char except "\"">* "\"" ;
IDENTIFIER     → ALPHA ( ALPHA | DIGIT )* ;
ALPHA          → "a" ... "z" | "A" ... "Z" | "_" ;
DIGIT          → "0" ... "9" ;
KEYWORDS       → "SELECT" | "FROM" | "WHERE" | "GROUP BY" | "OFFSET" | "LIMIT"
*/

func Scan(source string) ([]Token, error) {
	if len(source) == 0 {
		return []Token{}, nil
	}

	switch c := source[0:1]; c {
	case ",":
		return recurse(SimpleToken(COMMA, c), source[1:])
	case "'":
		return scanStr("", source[1:])
	case " ", "\n", "\t":
		return Scan(source[1:])
	case "S":
		switch {
		case check(source, "SELECT"):
			return recurse(SimpleToken(SELECT, source[:6]), source[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'S'")
		}

	default:
		return nil, fmt.Errorf("unexpected character '%s'.", c)
	}
}

func recurse(tk Token, source string) ([]Token, error) {
	switch tokens, err := Scan(source); err {
	case nil:
		return append([]Token{tk}, tokens...), nil
	default:
		return nil, err
	}
}

func scanStr(s string, source string) ([]Token, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("reached end of input parsing string.")
	}

	switch c := source[0:1]; c {
	case "'":
		return recurse(NewToken(STRING, s, s), source[1:])
	default:
		return scanStr(s+c, source[1:])
	}
}

func check(source string, keyword string) bool {
	if len(source) < len(keyword) {
		return false
	}
	// TODO check case insensitive
	return source[:len(keyword)] == keyword
}
