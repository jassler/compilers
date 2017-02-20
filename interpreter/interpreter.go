package interpreter

import (
	"fmt"

	"github.com/compilers/parser"
)

// Interpret interprets program from given expression
func Interpret(expr *parser.Expression) {

	res := interpretExpression(expr)
	fmt.Println(res)

}

// interpretExpression evaluates the type of Expression we have.
// Possible types of expression:
// - ExprInteger: Return int value of Expression
// - ExprIf:      Call interpretIf function that further evaluates the result of our given expression
func interpretExpression(expr *parser.Expression) int64 {
	switch (*expr).(type) {
	case parser.ExprInteger:
		return (*expr).(parser.ExprInteger).GetInteger()

	case parser.ExprIf:
		exprIf := (*expr).(parser.ExprIf)
		return interpretIf(&exprIf)

	default:
		panic("No valid expression received")
	}
}

// interpretIf checks, if condition is not 0. If that's the case, we return the int value of our consequent, else the alternative.
func interpretIf(expr *parser.ExprIf) int64 {
	if interpretExpression(expr.GetCondition()) != 0 {
		return interpretExpression(expr.GetConsequent())
	}
	return interpretExpression(expr.GetAlternative())
}
