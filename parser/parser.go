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
	err := list.New()

	e := parseExpression(&scan, err)

	if err.Len() == 0 {
		err = nil
	}
	return e, err
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

// parseExpression turns tokens into an expression
func parseExpression(scan *scanner.Scanner, list *list.List) *Expression {
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

	return &e
}

// parseIf tries to turn tokens into valid if expression.
// Semantic: if <expr> then <expr> else <expr>
// If semantic rule's not met, push error into our list
func parseIf(scan *scanner.Scanner, list *list.List) *ExprIf {
	var e ExprIf

	// condition comes from an expression
	e.condition = parseExpression(scan, list)
	if e.condition == nil {
		return nil
	}

	t, succ := scan.NextToken()

	// if not successful, then we've reached the end of our token array
	if !succ {
		list.PushBack(createError("Unexpected ending. Expected token: "+scanner.GetStringFromTokenID(scanner.TokenThen), t))
		return nil
	}

	// Expected token: then
	err := checkNextID(t.GetID(), scanner.TokenThen, t)

	if err != nil {
		list.PushBack(err)
		return nil
	}

	// consequent comes from an expression
	e.consequent = parseExpression(scan, list)
	if e.consequent == nil {
		return nil
	}

	// Expected token: else
	t, succ = scan.NextToken()
	if !succ {
		list.PushBack(createError("Unexpected ending. Expected token: "+scanner.GetStringFromTokenID(scanner.TokenElse), t))
		return nil
	}

	err = checkNextID(t.GetID(), scanner.TokenElse, t)
	if err != nil {
		list.PushBack(err)
		return nil
	}

	// alternative comes from an expression
	e.alternative = parseExpression(scan, list)
	if e.alternative == nil {
		return nil
	}

	// Expected token: end
	t, succ = scan.NextToken()
	if !succ {
		list.PushBack(createError("Unexpected ending. Expected token: "+scanner.GetStringFromTokenID(scanner.TokenEnd), t))
		return nil
	}

	err = checkNextID(t.GetID(), scanner.TokenEnd, t)
	if err != nil {
		list.PushBack(err)
	}

	return &e
}
