package parser

import "fmt"

// ExprIf has a condition. If condition is true, consequent will be executed, else alternativ.
type ExprIf struct {
	condition   Expression
	consequent  Expression
	alternative Expression
}

// Execute executes consequent condition, if condition.Execute() != 0, else execute alternative
func (i ExprIf) Execute() int64 {
	if i.condition.Execute() != 0 {
		return i.consequent.Execute()
	}

	return i.alternative.Execute()
}

// PrintExpr prints "if", then the expressions seperately
func (i ExprIf) PrintExpr(indent int) {
	printIndent(indent)
	fmt.Println("if")
	i.condition.PrintExpr(indent + 1)
	i.consequent.PrintExpr(indent + 1)
	i.alternative.PrintExpr(indent + 1)
}

// GetCondition returns condition of if expression
func (i *ExprIf) GetCondition() Expression {
	return i.condition
}

// GetConsequent returns consequent of if expression
func (i *ExprIf) GetConsequent() Expression {
	return i.consequent
}

// GetAlternative returns alternative of if expression
func (i *ExprIf) GetAlternative() Expression {
	return i.alternative
}
