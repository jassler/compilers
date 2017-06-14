package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"strconv"

	"github.com/simplang/interpreter"
	"github.com/simplang/lexer"
	"github.com/simplang/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: simplang <filename> [args]")
		return
	}

	path := os.Args[1]
	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("Could not read file:", err.Error())
		return
	}

	l := lexer.New(string(file))
	p := parser.New(l)
	a := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Generated", len(p.Errors()), "error(s):")
		for i, val := range p.Errors() {
			fmt.Printf("%d: %s\n", i+1, val)
		}
		return
	}

	params := make([]int64, len(os.Args)-2)
	for i := 2; i < len(os.Args); i++ {
		params[i-2], err = strconv.ParseInt(os.Args[i], 10, 64)
		if err != nil {
			fmt.Println(os.Args, "could not be converted to an integer")
			return
		}
	}
	//a.Print(0)
	fmt.Println(interpreter.Interprete(a, params))
}
