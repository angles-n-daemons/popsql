package parser

type tokenType int

const (
	// single character tokens
	NONE tokenType = iota
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

var keywordLookup = map[string]tokenType{
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

type token struct {
	Type    tokenType
	Lexeme  string
	Literal any
}

func simpleToken(ttype tokenType, lexeme string) *token {
	return newToken(ttype, lexeme, nil)
}

func newToken(ttype tokenType, lexeme string, literal any) *token {
	return &token{ttype, lexeme, literal}
}
