package parser

import "github.com/compilers/scanner"

// ParseFile parses file from given file
func ParseFile(path string) {
	scan := scanner.NewScanner(path)

	Parse(scan)
}

// Parse parses scanned file
func Parse(scan scanner.Scanner) {

	var e Expression

	e = parseExpression(&scan)

	e.PrintExpr(0)
}

func checkNextId(got int, expected int) {
	if got != expected {
		panic("Did not get what was expected")
	}
}

func parseExpression(scan *scanner.Scanner) Expression {
	var e Expression

	// ignoring errors for now and hoping for the best
	t, _ := scan.NextToken()

	switch t.GetID() {
	case scanner.TokenInteger:
		e = ExprInteger{i: int64(t.GetValue().(int))}

	case scanner.TokenIf:
		e = parseIf(scan)

	default:
		panic("Empty expression")
	}

	return e
}

func parseIf(scan *scanner.Scanner) ExprIf {
	var e ExprIf

	e.condition = parseExpression(scan)

	t, _ := scan.NextToken()
	checkNextId(t.GetID(), scanner.TokenThen)
	e.consequent = parseExpression(scan)

	t, _ = scan.NextToken()
	checkNextId(t.GetID(), scanner.TokenElse)
	e.alternative = parseExpression(scan)

	t, _ = scan.NextToken()
	checkNextId(t.GetID(), scanner.TokenEnd)

	return e
}
