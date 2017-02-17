package main

import (
	"fmt"
)

type Expression interface {
	PrintExpr(indent int)
}

type ExprInteger struct {
	i int64
}

type ExprIf struct {
	condition   Expression
	consequent  Expression
	alternative Expression
}

func printIndent(indent int) {
	for x := 0; x < indent; x++ {
		fmt.Print("    ")
	}
}

func (i ExprInteger) PrintExpr(indent int) {
	printIndent(indent)
	fmt.Println(i)
}

func (i ExprIf) PrintExpr(indent int) {
	printIndent(indent)
	fmt.Println("if")
	i.condition.PrintExpr(indent + 1)
	i.consequent.PrintExpr(indent + 1)
	i.alternative.PrintExpr(indent + 1)
}
