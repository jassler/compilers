package parser

import "fmt"

// ExprInteger only contains a number
type ExprInteger struct {
	i int64
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
