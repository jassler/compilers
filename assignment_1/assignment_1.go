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

		scan, err := scanner.NewScanner(fileName)

		if err != nil {
			s := "s"
			if err.Len() == 1 {
				s = ""
			}

			fmt.Printf("%d Syntax error%s found:\n", err.Len(), s)

			for e := err.Front(); e != nil; e = e.Next() {
				fmt.Println((e.Value.(error)).Error())
			}

			continue
		}

		expr, err := parser.Parse(scan)
		if err != nil {
			s := "s"
			if err.Len() == 1 {
				s = ""
			}

			fmt.Printf("%d Semantic error%s found:\n", err.Len(), s)
			for e := err.Front(); e != nil; e = e.Next() {
				fmt.Println((e.Value.(error)).Error())
			}

			continue
		}

		interpreter.Interpret(expr)
	}

}
