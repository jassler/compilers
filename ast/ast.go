package ast

import (
	"github.com/simplang/token"
)

// Expression is everything
type Expression interface {
	Print(indent int)
}

// Integer => int64
type Integer struct {
	Token token.Token
	Value int64
}

// Ident => identifier
type Ident struct {
	Token token.Token
	Name  string
}

// Program => function {function}
type Program struct {
	Token     token.Token
	Functions []*Function
}

// Function => "let" ident {ident}+ "in" expr "end"
type Function struct {
	Token  token.Token
	Name   *Ident
	Params []*Ident
	Body   Expression
}

// FunctionCall => ident arg {arg}
// arg => "(" expr ")"
type FunctionCall struct {
	Token  token.Token
	Name   string
	Params []Expression
}

// IfExpression => "if" expr "then" expr "else" expr "end"
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

// UnaryExpression => op expr
type UnaryExpression struct {
	Token    token.Token
	Operator token.TokenType
	Operand  Expression
}

// BinaryExpression => expr op expr
type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

// Binding => ident "=" expr
type Binding struct {
	Token token.Token
	Ident *Ident
	Expr  Expression
}

// LetExpression => "let" bindings "in" expr "end"
type LetExpression struct {
	Token    token.Token
	Bindings []*Binding
	Expr     Expression
}

// LoopExpression => "loop" bindings "in" expr "end"
type LoopExpression struct {
	Token    token.Token
	Bindings []*Binding
	Expr     Expression
}

// Recur => "recur" arg {arg}
// arg = "(" expr ")"
type Recur struct {
	Token token.Token
	Args  []Expression
}
