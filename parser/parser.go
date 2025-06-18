package parser

import (
    . "glox/util"
    . "glox/token"
    . "glox/ast"
    "fmt"
    "errors"
)

type Parser struct {
    tokens []Token
    curr int
}

type ParseError struct {
    token Token
    msg string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("%v - %v", e.token, e.msg)
}

// Parser "constructor"
func NewParser(tokens []Token) *Parser {
    return &Parser{tokens, 0}
}

// function to start parsing tokens
func (p *Parser) Parse() []Stmt {
    var ret []Stmt
    var pe *ParseError
    for !p.isAtEnd() {
        val, err := p.declaration()
        if errors.As(err, &pe) {
            return nil
        }
        ret = append(ret, val)
    }

    return ret
}

// RULE declaration: varDecl | statement
func (p *Parser) declaration() (Stmt, error) {
    if p.match(VAR) {
        ret, err := p.varDecl()
        if err != nil {
            p.synchronize()
            return nil, err
        }
        return ret, nil
    }

    ret, err := p.statement()
    if err != nil {
        p.synchronize()
        return nil, err
    }
    return ret, nil
}

// RULE statement: exprStmt | forStmt | ifStmt | printStmt | whileStmt | block
func (p *Parser) statement() (Stmt, error) {
    if p.match(FOR) {
        return p.forStmt()
    }
    if p.match(IF) {
        return p.ifStmt()
    }
    if p.match(PRINT) {
        return p.printStmt()
    }
    if p.match(WHILE) {
        return p.whileStmt()
    }
    if p.match(LEFT_BRACE) {
        val, err := p.block()
        if err != nil { return nil, err }
        return NewBlock(val), nil
    }

    return p.exprStmt()
}

// RULE forStmt: "for" "(" ( varDecl | exprStmt | ";" )
//               expression? ";" expression? ")" statement
func (p *Parser) forStmt() (Stmt, error) {
    _, err := p.consume(LEFT_PAREN, "Exprect '(' after for")
    if err != nil { return nil, err }

    var initializer Stmt
    if p.match(SEMICOLON) {
        initializer = nil
    } else if p.match(VAR) {
        initializer, err = p.varDecl()
        if err != nil { return nil, err }
    } else {
        initializer, err = p.exprStmt()
        if err != nil { return nil, err }
    }

    var condition Expr = nil
    if !p.check(SEMICOLON) {
        condition, err = p.expression()
        if err != nil { return nil, err }
    }

    _, err = p.consume(SEMICOLON, "Expect ';' after loop condition")
    if err != nil { return nil, err }

    var increment Expr = nil
    if !p.check(RIGHT_PAREN) {
        increment, err = p.expression()
        if err != nil { return nil, err }
    }

    _, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses")
    if err != nil { return nil, err }

    body, err := p.statement()
    if err != nil { return nil, err }

    if increment != nil {
        body = NewBlock( []Stmt{body, NewStmtExpression(increment)} )
    }

    if condition == nil {
        condition = NewLiteral(true)
    }
    body = NewWhile(condition, body)

    if initializer != nil {
        body = NewBlock( []Stmt{initializer, body} )
    }

    return body, nil
}

// RULE whileStmt: "while" "(" expression ")" statement
func (p *Parser) whileStmt() (Stmt, error) {
    _, err := p.consume(LEFT_PAREN, "Exprect '(' after while")
    if err != nil { return nil, err }

    condition, err := p.expression()
    if err != nil { return nil, err }

    _, err = p.consume(RIGHT_PAREN, "Expect ')' after condition")
    if err != nil { return nil, err }

    body, err := p.statement()
    if err != nil { return nil, err }

    return NewWhile(condition, body), nil
}

// RULE ifStmt: "if" "(" expression ")" statement ( "else" statement )?
func (p *Parser) ifStmt() (Stmt, error) {
    _, err := p.consume(LEFT_PAREN, "Expect '(' after if")
    if err != nil { return nil, err }

    condition, err := p.expression()
    if err != nil { return nil, err }

    _, err = p.consume(RIGHT_PAREN, "Expect ')' after condition")
    if err != nil { return nil, err }

    thenBranch, err := p.statement()
    if err != nil { return nil, err }

    var elseBranch Stmt = nil
    if p.match(ELSE) {
        elseBranch, err = p.statement()
        if err != nil { return nil, err }
    }

    return NewIf(condition, thenBranch, elseBranch), nil
}

// RULE block: "{" declaration* "}"
func (p *Parser) block() ([]Stmt, error) {
    var statements []Stmt

    for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
        val, err := p.declaration()
        if err != nil { return nil, err }
        statements = append(statements, val)
    }

    _, err := p.consume(RIGHT_BRACE, "Expect '}' after block")
    if err != nil { return nil, err }

    return statements, nil
}

// RULE printStmt: "print" expression ";"
func (p *Parser) printStmt() (Stmt, error) {
    val, err := p.expression()
    if err != nil { return nil, err }

    _, err = p.consume(SEMICOLON, "Expect ';' after value")
    if err != nil { return nil, err }

    return NewPrint(val), nil
}

// RULE "var" IDENTIFIER ( "=" expression )? ";"
func (p *Parser) varDecl() (Stmt, error) {
    name, err := p.consume(IDENTIFIER, "Expect variable name")
    if err != nil { return nil, err }

    var initializer Expr = nil
    if p.match(EQUAL) {
        initializer, err = p.expression()
        if err != nil { return nil, err }
    }

    _, err = p.consume(SEMICOLON, "Expect ';' after variable declaration")
    if err != nil { return nil, err }

    return NewVar(name, initializer), nil
}

// RULE exprStmt: expression ";"
func (p *Parser) exprStmt() (Stmt, error) {
    expr, err := p.expression()
    if err != nil { return nil, err }
    _, err = p.consume(SEMICOLON, "Expect ';' after value")

    if err != nil { return nil, err }

    return NewStmtExpression(expr), nil
}

// RULE expression: assignment
func (p *Parser) expression() (Expr, error) {
    return p.assignment()
}

// RULE assignment: IDENTIFIER "=" assignment | logic_or
func (p *Parser) assignment() (Expr, error) {
    expr, error := p.or()
    if error != nil { return nil, error }

    if p.match(EQUAL) {
        switch expr.(type) {
        case Variable:
            value, err := p.assignment()
            if err != nil { return nil, err }

            return NewAssign(expr.(Variable).Name, value), nil
        default:
            equals := p.previous()
            _, error := p.assignment()
            if error != nil { return nil, error }

            return nil, err(equals, "Invalid assignment target")
        }
    }

    return expr, nil 
}

// RULE logic_or: logic_and ( "or" logic_and )*
func (p *Parser) or() (Expr, error) {
    expr, err := p.and()
    if err != nil { return nil, err }

    for p.match(OR) {
        operator := p.previous()
        right, err := p.and()
        if err != nil { return nil, err }

        expr = NewLogical(expr, operator, right)
    }

    return expr, nil
}

// RULE logic_and: equality ( "and" equality )*
func (p *Parser) and() (Expr, error) {
    expr, err := p.equality()
    if err != nil { return nil, err }

    for p.match(OR) {
        operator := p.previous()
        right, err := p.equality()
        if err != nil { return nil, err }

        expr = NewLogical(expr, operator, right)
    }

    return expr, nil
}

// RULE equality: comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() (Expr, error) {
    expr, err := p.comparison()
    if err != nil { return nil, err }

    for p.match(BANG_EQUAL, EQUAL_EQUAL) {
        operator := p.previous()
        right, err := p.comparison()
        if err != nil { return nil, err }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE comparison: term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (Expr, error) {
    expr, err := p.term()
    if err != nil { return nil, err }

    for p.match(GREAT, GREAT_EQUAL, LESS, LESS_EQUAL) {
        operator := p.previous()
        right, err := p.term()
        if err != nil { return nil, err }

        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE term: factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() (Expr, error) {
    expr, err := p.factor()
    if err != nil { return nil, err }

    for p.match(MINUS, PLUS) {
        operator := p.previous()
        right, err := p.factor()
        if err != nil { return nil, err }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE factor: unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() (Expr, error) {
    expr, err := p.unary()
    if err != nil { return nil, err }

    for p.match(SLASH, STAR) {
        operator := p.previous()
        right, err := p.unary()
        if err != nil { return nil, err }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE unary: ( "!" | "-" ) unary | primary
func (p *Parser) unary() (Expr, error) {
    if p.match(BANG, MINUS) {
        operator := p.previous()
        right, err := p.unary()
        if err != nil { return nil, err }

        return NewUnary(operator, right), err
    }

    return p.primary()
}

// RULE primary: NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER
func (p *Parser) primary() (Expr, error) {
    switch {
    case p.match(FALSE):
        return NewLiteral(false), nil
    case p.match(TRUE):
        return NewLiteral(true), nil
    case p.match(NIL):
        return NewLiteral(nil), nil
    case p.match(NUMBER, STRING):
        return NewLiteral(p.previous().Literal), nil
    case p.match(IDENTIFIER):
        return NewVariable(p.previous()), nil
    case p.match(LEFT_PAREN):
        expr, err := p.expression()
        if err != nil { return nil, err }

        _, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
        if err != nil { return nil, err }

        return NewGrouping(expr), nil
    }

    return nil, err(p.peek(), "Expect expression.")
}

// function to check if the current token is any of the passed in types
// this does consume the token
func (p *Parser) match(types ...TokenType) bool {
    for _, ttype := range types {
        if (p.check(ttype)) {
            p.advance()
            return true
        }
    }

    return false
}

// function to consume current token if the type matches else return an error
func (p *Parser) consume(ttype TokenType, msg string) (Token, error) {
    if p.check(ttype) {
        return p.advance(), nil
    }

    return Token{}, err(p.peek(), msg)
}

// function to check if the given type matches the current token's type
func (p *Parser) check(ttype TokenType) bool {
    if p.isAtEnd() {
        return false
    }

    return p.peek().Type == ttype
}

// function to consume current token and move to the next one
func (p *Parser) advance() Token {
    if !p.isAtEnd() {
        p.curr++
    }

    return p.previous()
}

// function to check if at end of token list
func (p *Parser) isAtEnd() bool {
    return p.peek().Type == EOF
}

// function to return the token at the curr pointer without consuming
func (p *Parser) peek() Token {
    return p.tokens[p.curr]
}

// function to return the most recently consumed token (position = p.cur - 1)
func (p *Parser) previous() Token {
    if (p.curr == 0) {
        return Token{}
    }
    return p.tokens[p.curr - 1]
}

// return error
func err(token Token, msg string) error {
    TokenError(token, msg)
    return &ParseError{token, msg}
}

// function to go into panic mode and try to recover
func (p *Parser) synchronize() {
    p.advance()

    for !p.isAtEnd() {
        if p.previous().Type == SEMICOLON {
            return
        }

        switch p.peek().Type {
        case CLASS:
            fallthrough
        case FUN:
            fallthrough
        case VAR:
            fallthrough
        case FOR:
            fallthrough
        case IF:
            fallthrough
        case WHILE:
            fallthrough
        case PRINT:
            fallthrough
        case RETURN:
            return
        }

        p.advance()
    }
}
