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

	// data type
	DATATYPE_BOOLEAN
	DATATYPE_STRING
	DATATYPE_NUMBER

	// keywords
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

	"AND": AND,
	"OR":  OR,
	"NOT": NOT,

	"VALUES": VALUES,

	"NUMBER":  DATATYPE_NUMBER,
	"NUM":     DATATYPE_NUMBER,
	"INTEGER": DATATYPE_NUMBER,
	"INT":     DATATYPE_NUMBER,
	"CHAR":    DATATYPE_STRING,
	"VARCHAR": DATATYPE_STRING,
	"STRING":  DATATYPE_STRING,
	"BOOLEAN": DATATYPE_BOOLEAN,
	"BOOL":    DATATYPE_BOOLEAN,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

func (t *Token) Equal(o *Token) bool {
	return o != nil &&
		t.Type == o.Type &&
		t.Lexeme == o.Lexeme &&
		t.Literal == o.Literal
}

func simpleToken(ttype TokenType, lexeme string) *Token {
	return newToken(ttype, lexeme, nil)
}

func newToken(ttype TokenType, lexeme string, literal any) *Token {
	return &Token{ttype, lexeme, literal}
}
