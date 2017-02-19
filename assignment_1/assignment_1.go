package main

import "github.com/compilers/parser"

// main
func main() {
	fileName := "../src/github.com/compilers/sampleIF.txt"

	parser.ParseFile(fileName)
}
