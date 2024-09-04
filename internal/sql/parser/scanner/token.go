package scanner

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

func SimpleToken(ttype TokenType, lexeme string) *Token {
	return NewToken(ttype, lexeme, nil)
}

func NewToken(ttype TokenType, lexeme string, literal any) *Token {
	return &Token{ttype, lexeme, literal}
}
