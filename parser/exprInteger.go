package parser

import "fmt"

// ExprInteger only contains a number
type ExprInteger struct {
	i int64
}

// Execute returns i value
func (i ExprInteger) Execute() int64 {
	return i.i
}

// PrintExpr prints out the number from the struct
func (i ExprInteger) PrintExpr(indent int) {
	printIndent(indent)
	fmt.Println(i.i)
}

// GetInteger returns integer value of expression
func (i ExprInteger) GetInteger() int64 {
	return i.i
}
