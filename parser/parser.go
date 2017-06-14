package parser

import (
	"fmt"
	"strconv"

	"github.com/simplang/ast"
	"github.com/simplang/lexer"
	"github.com/simplang/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// read two tokens so curToken and peekToken is set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Token: p.curToken}

	if p.curToken.Type != token.LET {
		p.errors = append(p.errors, fmt.Sprintf("function has to start with 'let' (line %d.%d)", p.curToken.Line, p.curToken.Column))
		return nil
	}

	program.Functions = []*ast.Function{p.parseFunction()}

	for !p.peekTokenIs(token.EOF) {
		if !p.expectPeek(token.LET) {
			return nil
		}

		program.Functions = append(program.Functions, p.parseFunction())
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead (line %d.%d)", t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	var err error
	p.peekToken, err = p.l.NextToken()

	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
}

func (p *Parser) parseFunction() *ast.Function {
	f := &ast.Function{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	f.Name = &ast.Ident{Token: p.curToken, Name: p.curToken.Literal}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	f.Params = []*ast.Ident{&ast.Ident{Token: p.curToken, Name: p.curToken.Literal}}

	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
		f.Params = append(f.Params, &ast.Ident{Token: p.curToken, Name: p.curToken.Literal})
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	f.Body = p.parseExpression()

	if !p.expectPeek(token.END) {
		return nil
	}

	return f
}

// if precedence is not set, we know we're not inside a binary expression
func (p *Parser) parseExpression(precedence ...int) ast.Expression {
	var expr ast.Expression

	switch p.curToken.Type {
	case token.INT:
		expr = p.parseInteger()

	case token.IF:
		expr = p.parseIf()

	case token.NOT, token.MINUS:
		expr = p.parseUnary()

	case token.LPAREN:
		expr = p.parseLParen()

	case token.LET, token.LOOP:
		expr = p.parseLet()

	case token.IDENT:
		if p.peekTokenIs(token.LPAREN) {
			expr = p.parseFunctionCall()
		} else {
			expr = &ast.Ident{Token: p.curToken, Name: p.curToken.Literal}
		}

	case token.RECUR:
		expr = p.parseRecur()

	default:
		p.errors = append(p.errors, fmt.Sprintf("parser encountered an unexpected token type: %s (line %d.%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	if token.IsBinaryOperator(p.peekToken.Type) {

		// if we don't have a precedence set as a parameter, then we're at the top of the "calculation tree"
		// from there we need to continue eating all the operators until we're at the end
		if len(precedence) == 0 {
			for token.IsBinaryOperator(p.peekToken.Type) {
				p.nextToken()
				expr = p.parseBinaryOperator(expr)
			}
		} else if token.GetPrecedence(p.peekToken.Type) > precedence[0] {
			p.nextToken()
			expr = p.parseBinaryOperator(expr)
		}
	}

	return expr
}

// expr = integer
func (p *Parser) parseInteger() *ast.Integer {
	i := &ast.Integer{Token: p.curToken}
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)

	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Parser could not convert string to int. Value=%s, Error message=%s (line %d.%d)", p.curToken.Literal, err.Error(), p.curToken.Line, p.curToken.Column))
		return nil
	}

	i.Value = val
	return i
}

// expr = "if" expr "then" expr "else" expr "end"
func (p *Parser) parseIf() *ast.IfExpression {
	ifexpr := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	ifexpr.Condition = p.parseExpression()

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()
	ifexpr.Consequence = p.parseExpression()

	if !p.expectPeek(token.ELSE) {
		return nil
	}

	p.nextToken()
	ifexpr.Alternative = p.parseExpression()

	if !p.expectPeek(token.END) {
		return nil
	}

	return ifexpr
}

// expr = unop expr
// unop = "!" | "-"
func (p *Parser) parseUnary() *ast.UnaryExpression {
	unexpr := &ast.UnaryExpression{Token: p.curToken}

	unexpr.Operator = p.curToken.Type
	p.nextToken()

	unexpr.Operand = p.parseExpression()

	return unexpr
}

// expr = "(" expr binop expr ")"
// binop = "&&" | "||" | "<" | "==" | "+" | "*"
func (p *Parser) parseBinaryOperator(left ast.Expression) *ast.BinaryExpression {
	biexpr := &ast.BinaryExpression{Token: p.curToken, Operator: p.curToken.Type, Left: left}

	// expect binary op
	if !token.IsBinaryOperator(p.curToken.Type) {
		p.errors = append(p.errors, fmt.Sprintf("expected token to be a binary operator (+ * && || == <), instead got=%s (line %d.%d)", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	precedence := token.GetPrecedence(p.curToken.Type)

	p.nextToken()
	biexpr.Right = p.parseExpression(precedence)

	if token.IsBinaryOperator(p.curToken.Type) && token.GetPrecedence(p.curToken.Type) > precedence {
		biexpr = p.parseBinaryOperator(biexpr)
	}

	return biexpr
}

func (p *Parser) parseLParen() ast.Expression {
	p.nextToken()

	if p.curToken.Type == token.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("an empty set of parentheses is not valid (line %d.%d)", p.curToken.Line, p.curToken.Column))
		return nil
	}

	e := p.parseExpression()
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return e
}

func (p *Parser) parseFunctionCall() *ast.FunctionCall {
	fc := &ast.FunctionCall{Token: p.curToken, Name: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fc.Params = []ast.Expression{p.parseLParen()}

	for p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		fc.Params = append(fc.Params, p.parseLParen())
	}

	return fc
}

// parses let AND loop
func (p *Parser) parseLet() ast.Expression {
	t := p.curToken

	if t.Type != token.LET && t.Type != token.LOOP {
		p.errors = append(p.errors, fmt.Sprintf("Internal error. Called parseLet() with wrong token (%v)\n", t))
		return nil
	}

	bind := []*ast.Binding{}

	bind = append(bind, p.parseBinding())
	for p.peekToken.Type == token.AND {
		p.nextToken()
		bind = append(bind, p.parseBinding())
	}

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()
	e := p.parseExpression()

	p.expectPeek(token.END)

	if t.Type == token.LET {
		return &ast.LetExpression{Token: t, Bindings: bind, Expr: e}
	}
	if t.Type == token.LOOP {
		return &ast.LoopExpression{Token: t, Bindings: bind, Expr: e}
	}
	return nil
}

// ident "=" expr
func (p *Parser) parseBinding() *ast.Binding {
	b := &ast.Binding{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	b.Ident = &ast.Ident{Token: p.curToken, Name: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	b.Expr = p.parseExpression()
	return b
}

// "recur" arg {arg}
func (p *Parser) parseRecur() *ast.Recur {
	rec := &ast.Recur{Token: p.curToken, Args: []ast.Expression{}}

	for {
		if !p.expectPeek(token.LPAREN) {
			return nil
		}

		p.nextToken()
		rec.Args = append(rec.Args, p.parseExpression())

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		if !p.peekTokenIs(token.LPAREN) {
			break
		}
	}

	return rec
}
