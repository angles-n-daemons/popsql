package scanner

type TokenType int

const (
	// single character
	COMMA = iota

	// literals
	STRING

	SELECT
)

func (tt TokenType) String() string {
	return [...]string{"COMMA", "STRING", "SELECT"}[tt]
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

func SimpleToken(ttype TokenType, lexeme string) Token {
	return NewToken(ttype, lexeme, nil)
}

func NewToken(ttype TokenType, lexeme string, literal any) Token {
	return Token{ttype, lexeme, literal}
}
