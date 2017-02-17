package main

// main
func main() {
	// fileName := "../src/github.com/compilers/sample.txt"

	// ScanFile(fileName)

	i1 := ExprInteger{1}
	i2 := ExprInteger{2}
	i3 := ExprInteger{3}

	con := ExprIf{
		condition:   i1,
		consequent:  i2,
		alternative: i3,
	}

	con.PrintExpr(0)
}
