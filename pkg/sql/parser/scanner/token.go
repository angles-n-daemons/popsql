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
	DATATYPE_NUMBER
	DATATYPE_STRING
	DATATYPE_BOOLEAN

	// KEYWORDS
	SELECT
	INSERT
	INTO
	UPDATE
	DELETE

	CREATE
	TABLE

	FROM
	WHERE
	GROUP
	OFFSET
	ORDER
	LIMIT
	SET

	AND
	OR
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
	"INTO":   INTO,
	"UPDATE": UPDATE,
	"DELETE": DELETE,

	"CREATE": CREATE,
	"TABLE":  TABLE,

	"FROM":   FROM,
	"WHERE":  WHERE,
	"GROUP":  GROUP,
	"OFFSET": OFFSET,
	"LIMIT":  LIMIT,
	"SET":    SET,

	"NUMBER":  DATATYPE_NUMBER,
	"STRING":  DATATYPE_STRING,
	"BOOLEAN": DATATYPE_BOOLEAN,

	"AND": AND,
	"OR":  OR,
	"NOT": NOT,

	"VALUES": VALUES,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

func simpleToken(ttype TokenType, lexeme string) *Token {
	return newToken(ttype, lexeme, nil)
}

func newToken(ttype TokenType, lexeme string, literal any) *Token {
	return &Token{ttype, lexeme, literal}
}
