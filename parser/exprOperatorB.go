package parser

import (
	"fmt"

	"github.com/compilers/scanner"
)

// ExprOperatorBinary consists of two expressions connected with an operator (eg. "+", "-", "&&", ...)
type ExprOperatorBinary struct {
	op int
	e1 Expression
	e2 Expression

	f func(Expression, Expression) int64
}

// Execute returns result of instruction
func (e ExprOperatorBinary) Execute() int64 {
	// Since we can't assign default values we have to make sure that the function is not nil
	if e.f != nil {
		return e.f(e.e1, e.e2)
	}

	// Probably could have turned that method around a little bit, but since I assume
	// that it'll work first try, I want to avoid too many branches
	e.setOperator(e.op)
	return e.f(e.e1, e.e2)
}

// PrintExpr prints operator, then the expression
func (e ExprOperatorBinary) PrintExpr(indent int) {
	printIndent(indent)
	fmt.Println(scanner.GetStringFromTokenID(e.op))
	e.e1.PrintExpr(indent + 1)
	e.e2.PrintExpr(indent + 1)
}

// GetExpressions returns the 2 expressions from our binary operator struct
func (e *ExprOperatorBinary) GetExpressions() (Expression, Expression) {
	return e.e1, e.e2
}

// GetOperator returns token ID of our operator
func (e *ExprOperatorBinary) GetOperator() int {
	return e.op
}

// setOperator sets operator id and its corresponding function
func (e *ExprOperatorBinary) setOperator(op int) error {
	e.op = op

	switch op {
	/* && */
	case scanner.TokenLogicAnd:
		e.f = execOperatorAnd

	/* || */
	case scanner.TokenLogicOr:
		e.f = execOperatorOr

	/* < */
	case scanner.TokenLess:
		e.f = execOperatorLess

	/* == */
	case scanner.TokenEquals:
		e.f = execOperatorEquals

	/* + */
	case scanner.TokenPlus:
		e.f = execOperatorPlus

	/* * */
	case scanner.TokenTimes:
		e.f = execOperatorTimes

	default:
		e.f = execNil
		return fmt.Errorf("Operator could not be assigned to a function: Operator ID = %d, Value = %s", op, scanner.GetStringFromTokenID(op))
	}

	return nil
}

// binary operators: "&&" | "||" | "<" | "==" | "+" | "*"
func execOperatorAnd(e1 Expression, e2 Expression) int64 {
	if e1.Execute() != 0 {
		if e2.Execute() != 0 {
			return 1
		}
	}

	return 0
}

func execOperatorOr(e1 Expression, e2 Expression) int64 {
	if e1.Execute() == 0 {
		if e2.Execute() == 0 {
			return 0
		}
	}

	return 1
}

func execOperatorLess(e1 Expression, e2 Expression) int64 {
	if e1.Execute() < e2.Execute() {
		return 1
	}

	return 0
}

func execOperatorEquals(e1 Expression, e2 Expression) int64 {
	if e1.Execute() == e2.Execute() {
		return 1
	}

	return 0
}

func execOperatorPlus(e1 Expression, e2 Expression) int64 {
	return e1.Execute() + e2.Execute()
}

func execOperatorTimes(e1 Expression, e2 Expression) int64 {
	return e1.Execute() * e2.Execute()
}

func execNil(e1 Expression, e2 Expression) int64 {
	fmt.Println("Nil executed for binary operator between 2 expressions")
	return 0
}
