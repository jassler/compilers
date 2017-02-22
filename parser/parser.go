package parser

import (
	"container/list"
	"fmt"

	"github.com/compilers/scanner"
)

// ParseFile parses file from given file
func ParseFile(path string) (*Expression, *list.List) {
	scan, _ := scanner.NewScanner(path)

	return Parse(scan)
}

// Parse parses scanned file
func Parse(scan scanner.Scanner) (*Expression, *list.List) {
	var e Expression
	err := list.New()

	e = parseExpression(&scan, err)

	if err.Len() == 0 {
		err = nil
	}

	return &e, err
}

// Creates error message
func createError(msg string, t *scanner.Token) error {
	if t == nil {
		return fmt.Errorf("Semantic error: %s", msg)
	}
	return fmt.Errorf("Semantic error in line %d: %s", t.GetLineNumber(), msg)
}

// checkNextID checks, if the next token coincides with what's expected. Returns nil if that's the case, else return an error
func checkNextID(got int, expected int, t *scanner.Token) error {
	if got != expected {
		// We did not ge what was expeced. Return an error that explains what token we expected
		return createError("Unexpected token. Expected: "+scanner.GetStringFromTokenID(expected)+", Received: "+scanner.GetStringFromTokenID(got), t)
	}
	// Valid token, don't return an error
	return nil
}

// checkNextToken checks, if next token fullfills expectation. Returns true if that's the case
func checkNextToken(scan *scanner.Scanner, expected int, list *list.List) bool {

	t, succ := scan.NextToken()

	// End of file reached
	if !succ {
		list.PushBack(createError("Unexpected ending. Expected token: "+scanner.GetStringFromTokenID(expected), t))
		return false
	}

	// Expected token
	err := checkNextID(t.GetID(), expected, t)

	if err != nil {
		list.PushBack(err)
		return false
	}

	return true
}

// parseExpression turns tokens into an expression
func parseExpression(scan *scanner.Scanner, list *list.List) Expression {
	var e Expression

	t, succ := scan.NextToken()

	// if not successful, then we have reached the end of our token array
	if !succ {
		list.PushBack(createError("Unexpeced ending", t))
		return nil
	}

	// based on the token ID, convert expression into integer, if-statement, ...
	switch t.GetID() {
	case scanner.TokenInteger:
		e = ExprInteger{i: int64(t.GetValue().(int))}

	case scanner.TokenIf:
		e = parseIf(scan, list)

	default:
		list.PushBack(createError("Unexpected token: "+scanner.GetStringFromTokenID(t.GetID()), t))
		return nil
	}

	return e
}

// parseIf tries to turn tokens into valid if expression.
// Semantic: if <expr> then <expr> else <expr>
// If semantic rule's not met, push error into our list
func parseIf(scan *scanner.Scanner, list *list.List) ExprIf {
	var e ExprIf

	// condition comes from an expression
	e.condition = parseExpression(scan, list)
	if e.condition == nil {
		return e
	}

	// Expected token: then
	if !checkNextToken(scan, scanner.TokenThen, list) {
		return e
	}

	// consequent comes from an expression
	e.consequent = parseExpression(scan, list)
	if e.consequent == nil {
		return e
	}

	// Expected token: else
	if !checkNextToken(scan, scanner.TokenElse, list) {
		return e
	}

	// alternative comes from an expression
	e.alternative = parseExpression(scan, list)
	if e.alternative == nil {
		return e
	}

	// Expected token: end
	if !checkNextToken(scan, scanner.TokenEnd, list) {
		return e
	}

	return e
}
