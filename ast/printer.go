package ast

import "fmt"

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("  ")
	}
}

// Print Integer
// 4
func (i *Integer) Print(indent int) {
	printIndent(indent)
	fmt.Println(i.Value)
}

// Print Identifier
// a
func (i *Ident) Print(indent int) {
	printIndent(indent)
	fmt.Println(i.Name)
}

// Print Program
// function
//     ...
func (p *Program) Print(indent int) {
	for _, f := range p.Functions {
		f.Print(indent + 1)
	}
}

// Print FunctionCall
// add
//   3
//   2
func (fc *FunctionCall) Print(indent int) {
	printIndent(indent)
	fmt.Println(fc.Name)
	for _, arg := range fc.Params {
		arg.Print(indent + 1)
	}
}

// Print Function
// function
//     main
//       a
//       b
//   +
//     a
//     b
func (f *Function) Print(indent int) {
	printIndent(indent)
	fmt.Println("function")
	f.Name.Print(indent + 2)
	for _, p := range f.Params {
		p.Print(indent + 3)
	}

	f.Body.Print(indent + 1)
}

// Print if expression
// if
//   1 (condition)
//   2 (consequence)
//   3 (alternative)
func (ie *IfExpression) Print(indent int) {
	printIndent(indent)
	fmt.Println("if")
	ie.Condition.Print(indent + 1)
	ie.Consequence.Print(indent + 1)
	ie.Alternative.Print(indent + 1)
}

// Print unasy expression
// !
//   7
func (ue *UnaryExpression) Print(indent int) {
	printIndent(indent)
	fmt.Println(ue.Operator)
	ue.Operand.Print(indent + 1)
}

// Print binary expression
// +
//   3
//   *
//     7
//     4
func (be *BinaryExpression) Print(indent int) {
	printIndent(indent)
	fmt.Println(be.Operator)
	be.Left.Print(indent + 1)
	be.Right.Print(indent + 1)
}

// Print Binding
// a
//   7
func (b *Binding) Print(indent int) {
	b.Ident.Print(indent)
	b.Expr.Print(indent + 1)
}

// Print Let Expression
// eg. for 'let a = 7 in a end' =>
// let
//     a
//       7
//   a
func (le *LetExpression) Print(indent int) {
	printIndent(indent)
	fmt.Println("let")

	for _, b := range le.Bindings {
		b.Print(indent + 2)
	}

	le.Expr.Print(indent + 1)
}

// Print Loop expression
// loop
//     a
//       7
//   expression
func (le *LoopExpression) Print(indent int) {
	printIndent(indent)
	fmt.Println("loop")

	for _, b := range le.Bindings {
		b.Print(indent + 2)
	}

	le.Expr.Print(indent + 1)
}

// Print recur
// recur
//   1
//   2
func (r *Recur) Print(indent int) {
	printIndent(indent)
	fmt.Println("recur")
	for _, arg := range r.Args {
		arg.Print(indent + 1)
	}
}
