package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

// Step 3 Pratt Parser
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_           int = iota
	LOWEST          //
	EQUALS          // ==
	LESSGREATER     // < or >
	SUM             // +
	PRODUCT         // *
	PREFIX          // -x or !x
	CALL            // myFunc(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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
	inflixParseFns map[token.TokenType]infixParseFn
}

// registerPrefix adds a Prefix entry to the map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInflix adds nan Inflix entry to the map
func (p *Parser) registerInflix(tokenType token.TokenType, fn infixParseFn) {
	p.inflixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
	p := &Parser{
		l:              l,
		errors:         []string{},
		inflixParseFns: make(map[token.TokenType]infixParseFn),
		prefixParseFns: make(map[token.TokenType]prefixParseFn)}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerInflix(token.EQ, p.parseInfixExpression) // Infix Exprsessions
	p.registerInflix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInflix(token.LT, p.parseInfixExpression)
	p.registerInflix(token.GT, p.parseInfixExpression)
	p.registerInflix(token.PLUS, p.parseInfixExpression)
	p.registerInflix(token.MINUS, p.parseInfixExpression)
	p.registerInflix(token.SLASH, p.parseInfixExpression)
	p.registerInflix(token.ASTERISK, p.parseInfixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

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
		return p.parseExpressionStatement()
	}
}

// TODO
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	//	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST) // 0 precedence.
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}

// parseExpression
func (p *Parser) parseExpression(precedence int) ast.Expression {

	//	defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	left := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.inflixParseFns[p.peekToken.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}

	return left
}

func (p *Parser) parsePrefixExpression() ast.Expression {

	//	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// parseInfixExpression takes Left expression to constdruct an infix expression node with it. Then it assigns the precedence of the current token (operator of the infix expression) to the local var precedence.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	//	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	//fmt.Printf(" Operator: %s   Left: %q  Right: %q\n", expression.Operator, expression.Left.String(), expression.Right.String())
	return expression
}

// parseLetStatement parses a let statement inside parseStatment switch.
func (p *Parser) parseLetStatement() *ast.LetStatement {

	//	defer untrace(trace("parseLetStatment"))
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

	//	defer untrace(trace("parseReturnStatement"))
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	//      TODO skipping until we meet ; .

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	//	defer untrace(trace("parseIntegerLiteral"))
	literal := &ast.IntegerLiteral{Token: p.curToken}

	out, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("failed to parse %q to integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	literal.Value = out
	return literal
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// ERROR util
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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
