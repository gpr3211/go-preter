package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // current token
	peekToken token.Token // next token
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l}
	p.nextToken() // advance both current and peek
	p.nextToken()
	return p
}

// nextToken Advances the scanner to next token. Similar to peekchar, but with tokens
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
