package interpreter

import (
	"fmt"

	"reflect"

	"github.com/simplang/ast"
	"github.com/simplang/token"
)

// public map of functions with their names
var functions map[string]*ast.Function

func Interprete(expression ast.Expression, params []int64) int64 {
	prog, ok := expression.(*ast.Program)

	if !ok {
		fmt.Println("Error: Expression is not a program")
		return 0
	}

	functions = map[string]*ast.Function{}

	for _, f := range prog.Functions {
		if _, ok := functions[f.Name.Name]; ok {
			throwError(fmt.Sprintf("Error: Function '%s' is already defined", f.Name.Name), &f.Token)
		}
		functions[f.Name.Name] = f
	}

	val, ok := functions["main"]

	if !ok {
		throwError("Error: Function 'main' could not be found", nil)
	}

	return interpreteFunction(val, params)
}

func interpreteFunction(f *ast.Function, params []int64) int64 {
	l := len(params)
	if l != len(f.Params) {
		throwError(fmt.Sprintf("Error: Function called with wrong amount of arguments. expected=%d, got=%d", len(f.Params), l), &f.Token)
	}

	env := &environment{elements: make([]*element, l)}
	for i, val := range f.Params {
		env.elements[i] = &element{name: val.Name, value: params[i]}
	}

	res, isRec := interpreteExpr(f.Body, env)

	if isRec != nil {
		throwError(fmt.Sprintf("Error: recur appeared after function %s ended. Is a loop missing?", f.Name.Name), &f.Token)
	}

	return res
}

// functions return tuples
// int64 is the result we get
// []int64 are the argument values when recur is called
func interpreteExpr(expr ast.Expression, env *environment) (int64, []int64) {
	var res int64
	var rec []int64

	switch t := expr.(type) {
	case *ast.Integer:
		res = (*ast.Integer)(t).Value

	case *ast.IfExpression:
		res, rec = interpreteIf(t, env)

	case *ast.UnaryExpression:
		res, rec = interpreteUnop(t, env)

	case *ast.BinaryExpression:
		res, rec = interpreteBinop(t, env)

	case *ast.LetExpression:
		res, rec = interpreteLet(t, env)

	case *ast.Ident:
		res = env.getValue(&(*ast.Ident)(t).Token)

	case *ast.FunctionCall:
		fc := (*ast.FunctionCall)(t)
		f, ok := functions[fc.Name]
		if !ok {
			throwError(fmt.Sprintf("function '%s' is not defined", fc.Name), &fc.Token)
			break
		}

		res = interpreteFunction(f, evalArgs(fc.Params, &fc.Token, env))

	case *ast.LoopExpression:
		res, rec = interpreteLoop(t, env)

	case *ast.Recur:
		rec = evalArgs(t.Args, &t.Token, env)

	default:
		throwError(fmt.Sprintf("type is not valid in expression. got=%s", reflect.TypeOf(expr)), nil)
	}

	return res, rec
}

func evalArgs(expr []ast.Expression, t *token.Token, env *environment) []int64 {
	res := make([]int64, len(expr))
	var isRec []int64

	for i, val := range expr {
		res[i], isRec = interpreteExpr(val, env)

		if isRec != nil {
			throwError("recur may not appear inside an argument. Is a loop missing?", t)
		}
	}

	return res
}

func interpreteIf(expr *ast.IfExpression, env *environment) (int64, []int64) {
	res, isRec := interpreteExpr(expr.Condition, env)

	if isRec != nil {
		throwError("recur statement may not appear as a condition in an if statement. Is a loop missing?", &expr.Token)
	}

	if res != 0 {
		res, isRec = interpreteExpr(expr.Consequence, env)
		return res, isRec
	}

	res, isRec = interpreteExpr(expr.Alternative, env)
	return res, isRec
}

func interpreteUnop(expr *ast.UnaryExpression, env *environment) (int64, []int64) {
	res, isRec := interpreteExpr(expr.Operand, env)

	if isRec != nil {
		throwError("recur may not be used in connection with a unary operator. Is a loop missing?", &expr.Token)
	}

	switch expr.Operator {
	case token.NOT:
		if res != 0 {
			return 0, nil
		}
		return 1, nil

	case token.MINUS:
		return -res, nil

	default:
		throwError(fmt.Sprintf("invalid unary operator. Expected ! or -, got %s instead", expr.Operator), &expr.Token)
		return 0, nil
	}
}

func interpreteBinop(expr *ast.BinaryExpression, env *environment) (int64, []int64) {
	l, isRecl := interpreteExpr(expr.Left, env)
	r, isRecr := interpreteExpr(expr.Right, env)

	if (isRecl != nil) || (isRecr != nil) {
		throwError("recur may not be used with a binary operator. Is a loop missing?", &expr.Token)
	}

	switch expr.Operator {
	case token.LOG_AND:
		if l == 0 {
			return 0, nil
		}
		if r == 0 {
			return 0, nil
		}
		return 1, nil

	case token.LOG_OR:
		if l != 0 {
			return 1, nil
		}
		if r != 0 {
			return 1, nil
		}
		return 0, nil

	case token.LESS:
		if l < r {
			return 1, nil
		}
		return 0, nil

	case token.EQUAL:
		if l == r {
			return 1, nil
		}
		return 0, nil

	case token.PLUS:
		return l + r, nil

	case token.TIMES:
		return l * r, nil

	default:
		throwError(fmt.Sprintf("invalid binary operator. Expected &&, ||, <, ==, + or -, got %s instead", expr.Operator), &expr.Token)
		return 0, nil
	}
}

func interpreteLet(expr *ast.LetExpression, env *environment) (int64, []int64) {
	var res int64
	var isRec []int64

	for _, b := range expr.Bindings {
		res, isRec = interpreteExpr(b.Expr, env)

		if isRec != nil {
			throwError("recur may not appear as an argument. Is a loop missing?", &expr.Token)
		}
		env.appendElement(&element{name: b.Ident.Name, value: res})
	}

	res, isRec = interpreteExpr(expr.Expr, env)
	env.removeElements(len(expr.Bindings))

	return res, isRec
}

func interpreteLoop(expr *ast.LoopExpression, env *environment) (int64, []int64) {
	var res int64
	var isRec []int64

	for _, b := range expr.Bindings {
		res, isRec = interpreteExpr(b.Expr, env)

		if isRec != nil {
			throwError("recur may not appear as an argument. Is a loop missing?", &expr.Token)
		}

		env.appendElement(&element{name: b.Ident.Name, value: res})
	}

	for res, isRec = interpreteExpr(expr.Expr, env); isRec != nil; res, isRec = interpreteExpr(expr.Expr, env) {
		if len(isRec) != len(expr.Bindings) {
			throwError(fmt.Sprintf("recur has wrong amount of arguments. expected=%d, got=%d", len(expr.Bindings), len(isRec)), &expr.Token)
		}

		beg := len(env.elements) - len(expr.Bindings)
		for i, val := range isRec {
			env.elements[beg+i].value = val
		}
	}

	env.removeElements(len(expr.Bindings))

	return res, nil
}
