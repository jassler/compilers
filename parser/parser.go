package parser

import (
	"container/list"
	"fmt"

	"github.com/compilers/scanner"
)

// ParseFile parses file from given file
func ParseFile(path string) (Expression, *list.List) {
	scan, _ := scanner.NewScanner(path)

	return Parse(scan)
}

// Parse parses scanned file
func Parse(scan scanner.Scanner) (Expression, *list.List) {
	var e Expression
	err := list.New()

	e = parseExpression(&scan, err)

	if err.Len() == 0 {
		err = nil
	}

	return e, err
}

// Creates error message
func createError(msg string, t *scanner.Token) error {
	if t == nil {
		return fmt.Errorf("Syntax error: %s", msg)
	}
	return fmt.Errorf("Syntax error in line %d: %s", t.GetLineNumber(), msg)
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

	/* Integer */
	case scanner.TokenInteger:
		e = ExprInteger{i: int64(t.GetValue().(int))}

	/* if */
	case scanner.TokenIf:
		e = parseIf(scan, list)

	/* let */
	case scanner.TokenLet:
		e = parseLet(scan, list)

	/* unary operator ('!', '-') */
	case scanner.TokenNegate, scanner.TokenNot:
		e = ExprOperatorUnary{op: t.GetID(), expr: parseExpression(scan, list)}

	/* binary operator "(+, *, &&, ||, ==, <)" */
	case scanner.TokenOpenParen:
		e = parseBinaryOp(scan, list)

	default:
		list.PushBack(createError("Unexpected token: "+scanner.GetStringFromTokenID(t.GetID()), t))
		return nil
	}

	return e
}

// parseIf tries to turn tokens into valid if expression.
// Syntax: if <expr> then <expr> else <expr>
// If syntax rule's not met, push error into our list
func parseIf(scan *scanner.Scanner, list *list.List) Expression {
	var e ExprIf

	// condition comes from an expression
	e.condition = parseExpression(scan, list)
	if e.condition == nil {
		return nil
	}

	// Expected token: then
	if !checkNextToken(scan, scanner.TokenThen, list) {
		return nil
	}

	// consequent comes from an expression
	e.consequent = parseExpression(scan, list)
	if e.consequent == nil {
		return nil
	}

	// Expected token: else
	if !checkNextToken(scan, scanner.TokenElse, list) {
		return nil
	}

	// alternative comes from an expression
	e.alternative = parseExpression(scan, list)
	if e.alternative == nil {
		return nil
	}

	// Expected token: end
	checkNextToken(scan, scanner.TokenEnd, list)

	return e
}

func parseLet(scan *scanner.Scanner, list *list.List) Expression {
	var e = ExprLet{bindings: []Binding{}}

	// continue reading until token "in" comes up.
	// identifiers are connected with "and" token.
	for {
		// Expected token: identifier
		t, succ := scan.NextToken()

		if !succ {
			return nil
		}

		if t.GetID() != scanner.TokenIdentifier {
			list.PushBack(createError("Expected identifier after let.", t))
			return nil
		}

		b := Binding{ident: t.GetValue().(string)}

		// Expected token: =
		if !checkNextToken(scan, scanner.TokenEquals, list) {
			return nil
		}

		// Expeced token: expression
		if expr := parseExpression(scan, list); expr != nil {
			b.expr = &expr
		} else {
			return nil
		}

		// Expected token: "and" (continue) or "in" (break)
		t, succ = scan.NextToken()

		if !succ {
			return nil
		}

		if t.GetID() == scanner.TokenAnd {
			continue
		}

		if !checkNextToken(scan, scanner.TokenIn) {
			return nil
		}

		break
	}

	// Expected token: expression
	e.expr = parseExpression(scan, list)

	if e.expr == nil {
		return nil
	}

}

// Binary operator in the syntax: (a op b)
// op = "&&" , "||" , "<" , "==" , "+" , "*"
func parseBinaryOp(scan *scanner.Scanner, list *list.List) Expression {
	var e ExprOperatorBinary

	// a is another expression
	e.e1 = parseExpression(scan, list)
	if e.e1 == nil {
		return nil
	}

	// Expected token: Operator
	t, succ := scan.NextToken()

	if !succ {
		list.PushBack(createError("Unexpected ending. Expected an operator followed by an operand and a closing bracket", t))
		return nil
	}

	if err := e.setOperator(t.GetID()); err != nil {
		list.PushBack(createError(err.Error(), t))
	}

	// b is another expression
	e.e2 = parseExpression(scan, list)
	if e.e2 == nil {
		return nil
	}

	// Expected token: Closing parantheses ')'
	checkNextToken(scan, scanner.TokenCloseParen, list)

	return e
}
