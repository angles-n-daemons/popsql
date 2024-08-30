package scanner

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	// single character tokens
	NONE TokenType = iota
	COMMA
	LEFT_PAREN
	RIGHT_PAREN
	DOT
	MINUS
	PLUS
	STAR
	SLASH
	SEMICOLON

	// single or double character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// DATA_TYPES
	INTEGER
	VARCHAR
	FLOAT

	// KEYWORDS
	SELECT
	INSERT
	UPDATE
	DELETE

	CREATE

	FROM
	WHERE
	GROUP
	OFFSET
	LIMIT
	SET
	NOT

	VALUES
)

var keywordtrie = newTrie()

func init() {
	keywordtrie.seed(keywordLookup)
}

var keywordLookup = map[string]TokenType{
	",":      COMMA,
	"(":      LEFT_PAREN,
	")":      RIGHT_PAREN,
	".":      DOT,
	"-":      MINUS,
	"+":      PLUS,
	"*":      STAR,
	"/":      SLASH,
	";":      SEMICOLON,
	"!":      BANG,
	"!=":     BANG_EQUAL,
	"=":      EQUAL,
	"==":     EQUAL_EQUAL,
	">":      GREATER,
	">=":     GREATER_EQUAL,
	"<":      LESS,
	"<=":     LESS_EQUAL,
	"SELECT": SELECT,
	"INSERT": INSERT,
	"UPDATE": UPDATE,
	"DELETE": DELETE,

	"CREATE": CREATE,

	"FROM":   FROM,
	"WHERE":  WHERE,
	"GROUP":  GROUP,
	"OFFSET": OFFSET,
	"LIMIT":  LIMIT,
	"SET":    SET,
	"NOT":    NOT,

	"VALUES": VALUES,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

type trie struct {
	ttype    TokenType
	children map[byte]*trie
}

func newTrie() *trie {
	return &trie{
		children: map[byte]*trie{},
	}
}

func (t *trie) seed(tokens map[string]TokenType) {
	for s, ttype := range tokens {
		t.set(s, 0, ttype)
	}
}

func (t *trie) set(key string, i int, ttype TokenType) {
	if len(key) == i {
		fmt.Printf("setting type %s\n", ttype)
		fmt.Printf("setting type %s\n", ttype)
		t.ttype = ttype
		return
	}

	c := key[i]
	if t.children[c] == nil {
		t.children[c] = newTrie()
	}

	t.children[c].set(key, i+1, ttype)
}

func (t *trie) walk(i int) {
	if t.ttype != NONE {
		fmt.Println(strings.Repeat("-", i), i, t.ttype.String())
	}
	for c, child := range t.children {
		pad := ""
		if i > 1 {
			pad = strings.Repeat("-", i-1)
		}
		fmt.Println(pad, string(c))
		child.walk(i + 1)
	}
}

// conditions to look out for
//   - if we're not at the end and there's a value deeper, we use that value
//     when recursing, we add the character to search
//   - otherwise, if we have a type at this level, we return it
//   - otherwise we return nil
func (t *trie) get(search string, i int) *Token {
	// if the character has a child, look at the child
	if i < len(search) {
		if child, ok := t.children[search[i]]; ok {
			if token := child.get(search, i+1); token != nil {
				return token
			}
		}
	} else if t.ttype != NONE {
		token := SimpleToken(
			t.ttype,
			search[:i],
		)
		return &token
	}
	return nil
}

func SimpleToken(ttype TokenType, lexeme string) Token {
	return NewToken(ttype, lexeme, nil)
}

func NewToken(ttype TokenType, lexeme string, literal any) Token {
	return Token{ttype, lexeme, literal}
}
