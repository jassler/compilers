package token

/*
* Keywords are let, and, in, if, then, else, recur, loop, end.
* Operators are (, ), =, &&, ||, !, <, ==, +, *, -.
* Identifiers can contain only letters, digits, and the underscore, but cannot start with a digit.
* Integers are sequences of digits.
 */

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	// operators
	ASSIGN = "="

	LOG_AND = "&&"
	LOG_OR  = "||"
	NOT     = "!"
	LESS    = "<"
	EQUAL   = "=="

	PLUS  = "+"
	TIMES = "*"
	MINUS = "-"

	// delimiters
	LPAREN = "("
	RPAREN = ")"

	// keywords
	LET   = "let"
	AND   = "and"
	IN    = "in"
	IF    = "if"
	THEN  = "then"
	ELSE  = "else"
	RECUR = "recur"
	LOOP  = "loop"
	END   = "end"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"let":   LET,
	"and":   AND,
	"in":    IN,
	"if":    IF,
	"then":  THEN,
	"else":  ELSE,
	"recur": RECUR,
	"loop":  LOOP,
	"end":   END,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

const (
	PREC_LOWEST = iota
	PREC_LOGIC
	PREC_LEEQ
	PREC_PLUS
	PREC_TIMES
)

var binop = map[TokenType]int{
	LOG_AND: PREC_LOGIC,
	LOG_OR:  PREC_LOGIC,
	LESS:    PREC_LEEQ,
	EQUAL:   PREC_LEEQ,
	PLUS:    PREC_PLUS,
	TIMES:   PREC_TIMES,
}

func IsBinaryOperator(token TokenType) bool {
	_, ok := binop[token]
	return ok
}

func GetPrecedence(token TokenType) int {
	val, ok := binop[token]
	if !ok {
		return PREC_LOWEST
	}

	return val
}
