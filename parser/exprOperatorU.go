package parser

import (
	"fmt"

	"github.com/compilers/scanner"
)

// ExprOperatorUnary consists of either the unary character "!" or "-" and the following expression
type ExprOperatorUnary struct {
	op   int
	expr Expression
}

// Execute either returns -expr, !expr or just the expression if operator isn't valid
func (e ExprOperatorUnary) Execute() int64 {
	if e.op == scanner.TokenNegate {
		return -e.expr.Execute()
	}

	if e.op == scanner.TokenNot {
		return e.expr.Execute() ^ 1
	}

	return e.expr.Execute()
}

// PrintExpr prints operator, then the expression
func (e ExprOperatorUnary) PrintExpr(indent int) {
	printIndent(indent)
	fmt.Println(scanner.GetStringFromTokenID(e.op))
	e.expr.PrintExpr(indent + 1)
}

// GetOperator returns Token ID of unary operator set in ExprOperatorUnary struct
func (e *ExprOperatorUnary) GetOperator() int {
	return e.op
}

// GetExpression returns expression set after the unary operator
func (e *ExprOperatorUnary) GetExpression() Expression {
	return e.expr
}
