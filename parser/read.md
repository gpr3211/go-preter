### Parser

1. Pseudo-code
```go
    function parseProgram() {
    program = newProgramASTNode()
    advanceTokens()
    for (currentToken() != EOF_TOKEN) {
        statement = null
        if (currentToken() == LET_TOKEN) {
            statement = parseLetStatement()
        } else if (currentToken() == RETURN_TOKEN) {
            statement = parseReturnStatement()
        } else if (currentToken() == IF_TOKEN) {
        statement = parseIfStatement()
        }
        if (statement != null) {
        program.Statements.push(statement)
        }
        advanceTokens()
    }
    return program
}

    function parseLetStatement() {
        advanceTokens()
        identifier = parseIdentifier()
        advanceTokens()
        if currentToken() != EQUAL_TOKEN {
            parseError("no equal sign!")
            return null
        }
        advanceTokens()
        value = parseExpression()
        variableStatement = newVariableStatementASTNode()
        variableStatement.identifier = identifier
        variableStatement.value = value
        return variableStatement
    }
    function parseIdentifier() {
        identifier = newIdentifierASTNode()
        identifier.token = currentToken()
        return identifier
    }
    function parseExpression() {
        if (currentToken() == INTEGER_TOKEN) {
            if (nextToken() == PLUS_TOKEN)
            {
                return parseOperatorExpression()
            } else if (nextToken() == SEMICOLON_TOKEN) {
                return parseIntegerLiteral()
            }
        } else if (currentToken() == LEFT_PAREN) {
            return parseGroupedExpression()
        }
        // [...]
    }
    function parseOperatorExpression() {
        operatorExpression = newOperatorExpression()
        operatorExpression.left = parseIntegerLiteral()
        operatorExpression.operator = currentToken()
        operatorExpression.right = parseExpression()
        return operatorExpression()
    }
```
1. Since this is pseudocode there are a lot of omissions, of course. But the basic idea behind
    recursive-descent parsing is there. The entry point is parseProgram and it constructs the root
    node of the AST (newProgramASTNode()). It then builds the child nodes, the statements, by
    calling other functions that know which AST node to construct based on the current token.
2. These other functions call each other again, recursively.
        The most recursive part of this is in parseExpression and is only hinted at. But we can already
        see that in order to parse an expression like 5 + 5, we need to first parse 5 + and then call
        parseExpression() again to parse the rest, since after the + might be another operator expression,
like this: 5 + 5 * 10.
3. We will get to this later and look at expression parsing in detail, since it’s
    probably the most complicated but also the most beautiful part of the parser, making heavy
    use of “Pratt parsing”.
