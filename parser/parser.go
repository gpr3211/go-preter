package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

// Step 3 Pratt Parser
type (
	prefixParseFn func() ast.Expression
	inflixParseFn func(ast.Expression) ast.Expression
)

// Parser has 3 fields
//   - l *lexer.Lexer
//   - curToken token.Token
//   - peekToken token.Token
type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token // current token
	peekToken      token.Token // next token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	inflixParseFns map[token.TokenType]inflixParseFn
}

// registerPrefix adds a Prefix entry to the map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInflix adds nan Inflix entry to the map
func (p *Parser) registerInflix(tokenType token.TokenType, fn inflixParseFn) {
	p.inflixParseFns[tokenType] = fn
}

// STEP 2 WE ADD ERROR HANDLING

// ParseProgram constructs a root node and builds an AST.
//   - func(p *Parser) ParseProgram() *ast.Program
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken() // advance both current and peek
	p.nextToken()
	return p
}

// Errors returns the error array of strings.
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError appends an error msg to the errors array.
//   - input is token.TokenType
//   - output err msg is fmt.Sprintf("Expected token %s -- Got %s", t, p.peekToken.Type)
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected token %s -- Got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken Advances the scanner to next token. Similar to peekchar, but with tokens
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// parseStatement reads the curToken type and proceeds accordingly.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatement parses a let statement inside parseStatment switch.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO skipping the expressions until we encounter a semicolon
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	//      TODO skipping until we meet ; .

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// curTokenIs checks if the current token is a specified token.Type.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// curTokenIs checks if the peeked token is a specified token.Type.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
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
