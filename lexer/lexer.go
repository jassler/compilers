package lexer

import "github.com/simplang/token"

import "fmt"

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position in input
	line         int  // which line of code is it
	column       int  // which column of code is it
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// readChar sets position to 0 and readPosition to 1
	l.column = -1
	l.line = 1
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.line = -1
		l.column = -1
	} else {
		l.ch = l.input[l.readPosition]
		if l.ch != '\n' {
			l.column++
		} else {
			l.column = 0
			l.line++
		}
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() (token.Token, error) {
	var tok token.Token
	var err error

	l.skipWhitespace()

	startL := l.line
	startC := l.column

	// Operators are (, ), =, &&, ||, !, <, ==, +, *, -
	switch l.ch {
	case '=':
		if l.peekChar() != '=' {
			tok = newToken(token.ASSIGN, l.ch)
		} else {
			l.readChar()
			tok = token.Token{Type: token.EQUAL, Literal: "=="}
		}

	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '&':
		l.readChar()
		if l.ch == '&' {
			tok = token.Token{Type: token.LOG_AND, Literal: "&&"}
		} else {
			err = l.generateError(fmt.Sprint("Expected '&', got '", string(l.ch), "' instead"))
			tok = newToken(token.ILLEGAL, l.ch)
		}

	case '|':
		l.readChar()
		if l.ch == '|' {
			tok = token.Token{Type: token.LOG_OR, Literal: "||"}
		} else {
			err = l.generateError(fmt.Sprint("Expected '|', got '", string(l.ch), "' instead"))
			tok = newToken(token.ILLEGAL, l.ch)
		}

	case '!':
		tok = newToken(token.NOT, l.ch)
	case '<':
		tok = newToken(token.LESS, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '*':
		tok = newToken(token.TIMES, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readItentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			// we return here since we don't want to accidentally consume the next byte
			return tok, err
		}

		if isDigit(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok, err
		}

		err = l.generateError(fmt.Sprint("Invalid byte: got '", string(l.ch), "'"))
		tok = newToken(token.ILLEGAL, l.ch)

	}

	tok.Line = startL
	tok.Column = startC
	l.readChar()
	return tok, err
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch == '_')
}

func isDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9')
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || l.ch == '#' {
		// # for comments
		if l.ch == '#' {
			for l.ch != '\n' {
				l.readChar()
			}
		}

		l.readChar()

	}
}

func (l *Lexer) readItentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) generateError(msg string) error {
	return fmt.Errorf("Syntax error: %s (line %d.%d)", msg, l.line, l.column)
}
