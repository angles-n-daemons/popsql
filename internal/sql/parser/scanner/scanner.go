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

func Scan(s string) ([]Token, error) {
	if len(s) == 0 {
		return []Token{}, nil
	}

	switch c := s[0:1]; c {
	case ",":
		return recurse(SimpleToken(COMMA, s[:1]), Scan, s[1:])
	case "(":
		return recurse(SimpleToken(LEFT_PAREN, s[:1]), Scan, s[1:])
	case ")":
		return recurse(SimpleToken(RIGHT_PAREN, s[:1]), Scan, s[1:])
	case ".":
		return recurse(SimpleToken(DOT, s[:1]), Scan, s[1:])
	case "-":
		return recurse(SimpleToken(MINUS, s[:1]), Scan, s[1:])
	case "+":
		return recurse(SimpleToken(PLUS, s[:1]), Scan, s[1:])
	case "*":
		return recurse(SimpleToken(STAR, s[:1]), Scan, s[1:])
	case "/":
		return recurse(SimpleToken(SLASH, s[:1]), Scan, s[1:])
	case ";":
		return recurse(SimpleToken(SEMICOLON, s[:1]), Scan, s[1:])
	case "!":
		if len(s) > 1 && s[1:2] == "=" {
			return recurse(SimpleToken(BANG_EQUAL, s[:2]), Scan, s[2:])
		}
		return recurse(SimpleToken(BANG, s[:1]), Scan, s[1:])
	case "=":
		if len(s) > 1 && s[1:2] == "=" {
			return recurse(SimpleToken(EQUAL_EQUAL, s[:2]), Scan, s[2:])
		}
		return recurse(SimpleToken(EQUAL, s[:1]), Scan, s[1:])
	case ">":
		if len(s) > 1 && s[1:2] == "=" {
			return recurse(SimpleToken(GREATER_EQUAL, s[:2]), Scan, s[2:])
		}
		return recurse(SimpleToken(GREATER, s[:1]), Scan, s[1:])
	case "<":
		if len(s) > 1 && s[1:2] == "=" {
			return recurse(SimpleToken(LESS_EQUAL, s[:2]), Scan, s[2:])
		}
		return recurse(SimpleToken(LESS, s[:1]), Scan, s[1:])
	case "'":
		return scanStr("", s[1:])
	case "D":
		switch {
		case check(s, "DELETE"):
			return recurse(SimpleToken(DELETE, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'D'")
		}
	case "V":
		switch {
		case check(s, "VALUES"):
			return recurse(SimpleToken(VALUES, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'L'")
		}
	case "N":
		switch {
		case check(s, "NOT"):
			return recurse(SimpleToken(NOT, s[:3]), Scan, s[3:])
		default:
			return nil, fmt.Errorf("unmatched character 'L'")
		}
	case "L":
		switch {
		case check(s, "LIMIT"):
			return recurse(SimpleToken(LIMIT, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'L'")
		}
	case "O":
		switch {
		case check(s, "OFFSET"):
			return recurse(SimpleToken(OFFSET, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'O'")
		}
	case "G":
		switch {
		case check(s, "GROUP"):
			return recurse(SimpleToken(GROUP, s[:5]), Scan, s[5:])
		default:
			return nil, fmt.Errorf("unmatched character 'G'")
		}
	case "W":
		switch {
		case check(s, "WHERE"):
			return recurse(SimpleToken(WHERE, s[:5]), Scan, s[5:])
		default:
			return nil, fmt.Errorf("unmatched character 'W'")
		}
	case "C":
		switch {
		case check(s, "CREATE"):
			return recurse(SimpleToken(CREATE, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'C'")
		}
	case "F":
		switch {
		case check(s, "FROM"):
			return recurse(SimpleToken(FROM, s[:4]), Scan, s[4:])
		default:
			return nil, fmt.Errorf("unmatched character 'F'")
		}
	case "I":
		switch {
		case check(s, "INSERT"):
			return recurse(SimpleToken(INSERT, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'I'")
		}
	case "S":
		switch {
		case check(s, "SELECT"):
			return recurse(SimpleToken(SELECT, s[:6]), Scan, s[6:])
		case check(s, "SET"):
			return recurse(SimpleToken(SET, s[:3]), Scan, s[3:])
		default:
			return nil, fmt.Errorf("unmatched character 'S'")
		}
	case "U":
		switch {
		case check(s, "UPDATE"):
			return recurse(SimpleToken(UPDATE, s[:6]), Scan, s[6:])
		default:
			return nil, fmt.Errorf("unmatched character 'U'")
		}
	case " ", "\n", "\t":
		return Scan(s[1:])
	default:
		return nil, fmt.Errorf("unexpected character '%s'.", c)
	}
}

func ScanSequential(s string) ([]Token, error) {
	checkSeq := func(source string, i int, keyword string) bool {
		if len(source) < i+len(keyword) {
			return false
		}
		// TODO check case insensitive
		return source[:len(keyword)] == keyword
	}
	i := 0
	tokens := []Token{}
	for i < len(s) {
		switch c := s[0:1]; c {
		case ",":
			i++
			tokens = append(tokens, SimpleToken(COMMA, s[:1]))
		case "(":
			i++
			tokens = append(tokens, SimpleToken(LEFT_PAREN, s[:1]))
		case ")":
			i++
			tokens = append(tokens, SimpleToken(RIGHT_PAREN, s[:1]))
		case ".":
			i++
			tokens = append(tokens, SimpleToken(DOT, s[:1]))
		case "-":
			i++
			tokens = append(tokens, SimpleToken(MINUS, s[:1]))
		case "+":
			i++
			tokens = append(tokens, SimpleToken(PLUS, s[:1]))
		case "*":
			i++
			tokens = append(tokens, SimpleToken(STAR, s[:1]))
		case "/":
			i++
			tokens = append(tokens, SimpleToken(SLASH, s[:1]))
		case ";":
			i++
			tokens = append(tokens, SimpleToken(SEMICOLON, s[:1]))
		case "!":
			if len(s) > i+1 && s[1:2] == "=" {
				i += 2
				tokens = append(tokens, SimpleToken(BANG_EQUAL, s[:2]))
			} else {
				i++
				tokens = append(tokens, SimpleToken(BANG, s[:1]))
			}
		case "=":
			if len(s) > i+1 && s[1:2] == "=" {
				i += 2
				tokens = append(tokens, SimpleToken(EQUAL_EQUAL, s[:2]))
			} else {
				i++
				tokens = append(tokens, SimpleToken(EQUAL, s[:1]))
			}
		case ">":
			if len(s) > i+1 && s[1:2] == "=" {
				i += 2
				tokens = append(tokens, SimpleToken(GREATER_EQUAL, s[:2]))
			} else {
				i++
				tokens = append(tokens, SimpleToken(GREATER, s[:1]))
			}
		case "<":
			if len(s) > i+1 && s[1:2] == "=" {
				i += 2
				tokens = append(tokens, SimpleToken(LESS_EQUAL, s[:2]))
			} else {
				i++
				tokens = append(tokens, SimpleToken(LESS, s[:1]))
			}
		case "D":
			switch {
			case checkSeq(s, i, "DELETE"):
				i += 6
				tokens = append(tokens, SimpleToken(DELETE, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'D'")
			}
		case "V":
			switch {
			case checkSeq(s, i, "VALUES"):
				i += 6
				tokens = append(tokens, SimpleToken(VALUES, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'L'")
			}
		case "N":
			switch {
			case checkSeq(s, i, "NOT"):
				i += 3
				tokens = append(tokens, SimpleToken(NOT, s[:3]))
			default:
				return nil, fmt.Errorf("unmatched character 'L'")
			}
		case "L":
			switch {
			case checkSeq(s, i, "LIMIT"):
				i += 5
				tokens = append(tokens, SimpleToken(LIMIT, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'L'")
			}
		case "O":
			switch {
			case checkSeq(s, i, "OFFSET"):
				i += 6
				tokens = append(tokens, SimpleToken(OFFSET, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'O'")
			}
		case "G":
			switch {
			case checkSeq(s, i, "GROUP"):
				i += 5
				tokens = append(tokens, SimpleToken(GROUP, s[:5]))
			default:
				return nil, fmt.Errorf("unmatched character 'G'")
			}
		case "W":
			switch {
			case checkSeq(s, i, "WHERE"):
				i += 5
				tokens = append(tokens, SimpleToken(WHERE, s[:5]))
			default:
				return nil, fmt.Errorf("unmatched character 'W'")
			}
		case "C":
			switch {
			case checkSeq(s, i, "CREATE"):
				tokens = append(tokens, SimpleToken(CREATE, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'C'")
			}
		case "F":
			switch {
			case checkSeq(s, i, "FROM"):
				tokens = append(tokens, SimpleToken(FROM, s[:4]))
			default:
				return nil, fmt.Errorf("unmatched character 'F'")
			}
		case "I":
			switch {
			case checkSeq(s, i, "INSERT"):
				tokens = append(tokens, SimpleToken(INSERT, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'I'")
			}
		case "S":
			switch {
			case checkSeq(s, i, "SELECT"):
				tokens = append(tokens, SimpleToken(SELECT, s[:6]))
			case checkSeq(s, i, "SET"):
				tokens = append(tokens, SimpleToken(SET, s[:3]))
			default:
				return nil, fmt.Errorf("unmatched character 'S'")
			}
		case "U":
			switch {
			case check(s, "UPDATE"):
				tokens = append(tokens, SimpleToken(UPDATE, s[:6]))
			default:
				return nil, fmt.Errorf("unmatched character 'U'")
			}
		case " ", "\n", "\t":
			i++
			return Scan(s[1:])
		default:
			return nil, fmt.Errorf("unexpected character '%s'.", c)
		}
	}
	return tokens, nil
}

func ScanWithMap(s string) ([]Token, error) {
	if len(s) == 0 {
		return []Token{}, nil
	}
	switch s[:1] {
	case " ", "\n", "\t":
		return ScanWithMap(s[1:])
	default:
		for i := 6; i >= 1; i-- {
			if len(s) < i {
				continue
			}
			if ttype, ok := keywordLookup[s[:i]]; ok {
				return recurse(SimpleToken(ttype, s[:i]), ScanWithMap, s[i:])
			}
		}
	}
	return nil, fmt.Errorf("unexpected character '%s'.", s[:1])
}

func ScanWithTrie(s string) ([]Token, error) {
	if len(s) == 0 {
		return []Token{}, nil
	}
	switch s[:1] {
	case " ", "\n", "\t":
		return ScanWithTrie(s[1:])
	default:
		if token := keywordtrie.get(s, 0); token != nil {
			fwd := len(token.Lexeme)
			return recurse(*token, ScanWithTrie, s[fwd:])
		} else {
			return nil, fmt.Errorf("Unseen character '%s'", s[:1])
		}

	}
}

func recurse(tk Token, f func(string) ([]Token, error), source string) ([]Token, error) {
	switch tokens, err := f(source); err {
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
		return recurse(NewToken(STRING, s, s), Scan, source[1:])
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
