package main

import (
	"fmt"
	"os"

	"container/list"

	"github.com/compilers/interpreter"
	"github.com/compilers/parser"
	"github.com/compilers/scanner"
)

// printErrors returns true, if errors exist
func printErrors(l *list.List) bool {
	if l == nil {
		return false
	}

	if l.Len() == 0 {
		return false
	}

	ext := ""
	if l.Len() > 1 {
		ext = "s"
	}

	fmt.Printf("Found %d error%s\n", l.Len(), ext)

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println((e.Value.(error)).Error())
	}
	return true
}

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

		scan, errs := scanner.NewScanner(fileName)
		if printErrors(errs) {
			continue
		}

		expr, errs := parser.Parse(scan)
		if printErrors(errs) {
			continue
		}

		fmt.Println(interpreter.Interpret(expr))
	}

}
