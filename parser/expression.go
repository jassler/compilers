package parser

import (
	"fmt"
)

// Expression is an interface for every other Expression to import from
type Expression interface {
	PrintExpr(indent int)
}

func printIndent(indent int) {
	for x := 0; x < indent; x++ {
		fmt.Print("  ")
	}
}
