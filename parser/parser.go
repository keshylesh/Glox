package parser

import (
    . "glox/util"
    . "glox/token"
    . "glox/ast"
    "errors"
    "fmt"
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

func (p *Parser) Parse() Expr {
    ret, err := p.expression()
    var pe *ParseError
    if errors.As(err, &pe) {
        return nil
    } else if err == nil {
        return ret
    }
    return nil
}

// RULE expression: equality
func (p *Parser) expression() (Expr, error) {
    return p.equality()
}

// RULE equality: comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() (Expr, error) {
    expr, err := p.comparison()
    if err != nil {
        return nil, err
    }

    for p.match(BANG_EQUAL, EQUAL_EQUAL) {
        operator := p.previous()
        right, err := p.comparison()
        if err != nil {
            return nil, err
        }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE comparison: term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (Expr, error) {
    expr, err := p.term()
    if err != nil {
        return nil, err
    }

    for p.match(GREAT, GREAT_EQUAL, LESS, LESS_EQUAL) {
        operator := p.previous()
        right, err := p.term()
        if err != nil {
            return nil, err
        }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE term: factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() (Expr, error) {
    expr, err := p.factor()
    if err != nil {
        return nil, err
    }

    for p.match(MINUS, PLUS) {
        operator := p.previous()
        right, err := p.factor()
        if err != nil {
            return nil, err
        }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE factor: unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() (Expr, error) {
    expr, err := p.unary()
    if err != nil {
        return nil, err
    }

    for p.match(SLASH, STAR) {
        operator := p.previous()
        right, err := p.unary()
        if err != nil {
            return nil, err
        }
        expr = NewBinary(expr, operator, right)
    }

    return expr, nil
}

// RULE unary: ( "!" | "-" ) unary | primary
func (p *Parser) unary() (Expr, error) {
    if p.match(BANG, MINUS) {
        operator := p.previous()
        right, err := p.unary()
        if err != nil {
            return nil, err 
        }
        return NewUnary(operator, right), err
    }

    return p.primary()
}

// RULE primary: NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
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
    case p.match(LEFT_PAREN):
        expr, err := p.expression()
        if err != nil {
            return nil, err
        }
        _, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
        if err != nil {
            return nil, err
        }
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

func err(token Token, msg string) error {
    TokenError(token, msg)
    return &ParseError{token, msg}
}

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
