package main

const (
	TokenInteger = iota
	TokenIdentifier

	TokenAnd
	TokenElse
	TokenEnd

	// TokenIf usage: if <expr> then <expr> else <expr>
	TokenIf
	TokenIn
	TokenLet
	TokenLoop
	TokenRecur
	TokenThen

	TokenOpenParen
	TokenCloseParen
	TokenNot
	TokenLess
	TokenPlus
	TokenTimes
	TokenNegate

	TokenAssign

	TokenEquals
	TokenLogicAnd
	TokenLogicOr
)
