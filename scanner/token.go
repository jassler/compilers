package scanner

import (
	"reflect"
	"strconv"
)

const (
	TokenInteger = iota
	TokenIdentifier

	TokenAnd
	TokenElse
	TokenEnd

	// TokenIf usage: if <expr> then <expr> else <expr>
	TokenIf
	TokenIn
	TokenLet
	TokenLoop
	TokenRecur
	TokenThen

	TokenOpenParen
	TokenCloseParen
	TokenNot
	TokenLess
	TokenPlus
	TokenTimes
	TokenNegate

	TokenAssign

	TokenEquals
	TokenLogicAnd
	TokenLogicOr
)

var mapOfKeywords = map[string]int{
	"and":  TokenAnd,
	"else": TokenElse,
	"end":  TokenEnd,

	// TokenIf usage: if <expr> then <expr> else <expr> end
	"if":    TokenIf,
	"in":    TokenIn,
	"let":   TokenLet,
	"loop":  TokenLoop,
	"recur": TokenRecur,
	"then":  TokenThen,
}

var mapOfOperators = map[string]int{
	"(": TokenOpenParen,
	")": TokenCloseParen,
	"!": TokenNot,
	"<": TokenLess,
	"+": TokenPlus,
	"*": TokenTimes,
	"-": TokenNegate,

	"=": TokenAssign,

	"==": TokenEquals,
	"&&": TokenLogicAnd,
	"||": TokenLogicOr,
}

// Token describes the token type and value
type Token struct {
	tokenID  int
	tokenVal interface{}
	lineNum  int
}

// GetID returns token id of token
func (t Token) GetID() int {
	return t.tokenID
}

// GetValue returns value of token
func (t Token) GetValue() interface{} {
	return t.tokenVal
}

// GetLineNumber returns line number where the token appears in the source code
func (t Token) GetLineNumber() int {
	return t.lineNum
}

func mapString(s string, m map[string]int) int {
	val, ok := m[s]

	if !ok {
		return -1
	}

	return val
}

// MapKeyword assigns keyword to corresponding token id
func MapKeyword(s string) int {
	return mapString(s, mapOfKeywords)
}

// MapOperator assigns operator to corresponding token id
func MapOperator(s string) int {
	return mapString(s, mapOfOperators)
}

func makeArrayOfKeys(m map[string]int) []string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, len(keys))

	for index := 0; index < len(keys); index++ {
		strkeys[index] = keys[index].String()
	}

	return strkeys
}

// GetKeywords returns array of keywords in simplang
func GetKeywords() []string {
	return makeArrayOfKeys(mapOfKeywords)
}

// GetOperators returns array of operators in simplang
func GetOperators() []string {
	return makeArrayOfKeys(mapOfOperators)
}

// GetStringFromTokenID returns keyword string from given id
func GetStringFromTokenID(id int) string {
	var key string
	var val int

	for key, val = range mapOfKeywords {
		if val == id {
			return key
		}
	}

	for key, val = range mapOfOperators {
		if val == id {
			return key
		}
	}

	if id == TokenIdentifier {
		return "identifier"
	}

	if id == TokenInteger {
		return "integer"
	}

	return "invalid id (" + strconv.Itoa(id) + ")"
}
