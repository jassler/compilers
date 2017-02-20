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

func interpretIf(expr *parser.ExprIf) int64 {
	if interpretExpression(expr.GetCondition()) != 0 {
		return interpretExpression(expr.GetConsequent())
	}
	return interpretExpression(expr.GetAlternative())
}
