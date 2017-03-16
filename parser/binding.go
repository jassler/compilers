package parser

// Binding connects identifier name and its expression
type Binding struct {
	ident string
	expr  *Expression
}

// GetIdentifier returns identifier name of binding
func (b *Binding) GetIdentifier() string {
	return b.ident
}

// GetExpression returns expression of given binding
func (b *Binding) GetExpression() *Expression {
	return b.expr
}
