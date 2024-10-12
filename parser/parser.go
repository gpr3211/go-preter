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
	token.LPAREN:   CALL,
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
	p.registerPrefix(token.IF, p.parseIfExpression) // IF
	p.registerPrefix(token.TRUE, p.parseBoolean)    // bools
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

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
//
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
		fmt.Println("err no prefix found :", prefix)
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
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// IF
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block

}

// ERROR util
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// curTokenIs checks if the current token is a specified token.Type.
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}
