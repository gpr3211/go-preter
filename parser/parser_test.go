package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x =  5;
let y = 10;
let foobar = 835;
`
	l := lexer.New(input)
	p := New(l) // parser

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Wrong number of statements -- Expected 3 got %v", len(program.Statements))
	}
	checkParserErrors(t, p)

	tests := []struct {
		expectedID string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedID) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors() // exctract error array
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let'. Received = %T", s)
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got %T", s)
		return false
	}

	// test name
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value -- Expected %s -- Received %s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name -- Expected %s -- Received %s", name, letStmt.Name)
		return false
	}
	return true
}

// TESTING RETURN

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return test;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements not  len 3.")
	}
	for _, stmt := range program.Statements {
		returnStm, ok := stmt.(*ast.ReturnStatement) // type assertion
		if !ok {
			t.Errorf("statement not return -- have %T", stmt)
			continue
		}
		if returnStm.TokenLiteral() != "return" {
			t.Errorf("statement not return -- have %T", returnStm.TokenLiteral())
		}
	}
}
