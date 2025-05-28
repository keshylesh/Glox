package util

import "fmt"

// Enum to represent the tokens that can be scanned
type TokenType int

const (
    // NO_TYPE: Special token; used to make sure the keyword
    // mapping doesn't accidentally use an existing token
    // if a mapping doesn't exist it would return 0 which would
    // correspond to LEFT_PAREN if NO_TYPE were not there
    NO_TYPE TokenType = iota 
    
    // Single-byte tokens
    LEFT_PAREN
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    COMMA
    DOT
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR

    // One or two byte tokens
    BANG
    BANG_EQUAL
    EQUAL
    EQUAL_EQUAL
    GREAT
    GREAT_EQUAL
    LESS
    LESS_EQUAL

    // Literals
    IDENTIFIER
    STRING
    NUMBER

    // Keywords
    AND
    CLASS
    ELSE
    FALSE
    FUN
    FOR
    IF
    NIL
    OR
    PRINT
    RETURN
    SUPER
    THIS
    TRUE
    VAR
    WHILE

    // End Of File
    EOF
)

// following line to use the stringer tool to automatically generate
// string mappings for all the values of TokenType for human
// readability
//go:generate stringer -type=TokenType

// create a new Object type for convenience
type Object interface{}

// keywords to recognise and map tokens to
var Keywords = map[string]TokenType{
    "and": AND,
    "class": CLASS,
    "else": ELSE,
    "false": FALSE,
    "for": FOR,
    "fun": FUN,
    "if": IF,
    "nil": NIL,
    "or": OR,
    "print": PRINT,
    "return": RETURN,
    "super": SUPER,
    "this": THIS,
    "true": TRUE,
    "var": VAR,
    "while": WHILE,
}

type ParseError struct {
    token Token
    msg string
}

func (e *ParseError) Error() {
    fmt.Printf("%v - %v", e.token, e.msg)
}
