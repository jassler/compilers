package parser

// ExprInteger only contains a number
type ExprLet struct {
	bindings []Binding
	expr     Expression
}

// Execute returns i value
func (l ExprLet) Execute() int64 {
	// todo
}

// PrintExpr prints out the number from the struct
func (l ExprLet) PrintExpr(indent int) {
	// todo
}

// GetInteger returns integer value of expression
func (l ExprLet) GetInteger() int64 {
	// todo
}
