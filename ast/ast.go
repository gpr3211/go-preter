package ast

import "interpreter/token"

// Node general node interface
type Node interface {
	TokenLiteral() string
}

// Statement Statement Node interface
type Statement interface {
	Node
	statementNode()
}

// Expression node interface
type Expression interface {
	Node
	expressionNode()
}

// Program node is going to be the root of each AST
type Program struct {
	Statements []Statement
}

// interface method
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LET
type LetStatement struct {
	Token token.Token // LET token
	Name  *Identifier
	Value Expression
}

func (i *LetStatement) statementNode()       {}
func (i *LetStatement) TokenLiteral() string { return i.Token.Literal }

// RETURN
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {} // empty, just to satisfy interface
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type Identifier struct {
	Token token.Token // the token IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
