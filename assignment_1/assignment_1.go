package main

import (
	"fmt"
	"os"

	"github.com/compilers/interpreter"
	"github.com/compilers/parser"
	"github.com/compilers/scanner"
)

// main
func main() {
	lenArgs := len(os.Args)

	if lenArgs < 2 {
		fmt.Println("Usage: program <filename>")
		return
	}

	for _, fileName := range os.Args[1:] {
		if lenArgs > 2 {
			fmt.Println(fileName)
		}
		scan := scanner.NewScanner(fileName)
		expr := parser.Parse(scan)
		interpreter.Interpret(expr)
	}

}
