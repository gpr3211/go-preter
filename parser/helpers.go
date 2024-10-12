package parser

import (
	"fmt"
	"interpreter/token"
)

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// curTokenIs checks if the peeked token is a specified token.Type.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected token %s -- Got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken Advances the scanner to next token. Similar to peekchar, but with tokens
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// expectPeek assertion functions. Enforces correctness of the token order by checking type of next token.
// Check type of peek and only advances if it is the correct.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekPrecedence peeks the next token. If empty returns 0.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence returns current precedence from table
// - It tells that +( token.Plus) and - have the same precedence
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
