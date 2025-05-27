package token

import (
    . "glox/util"
)

type Token struct {
    Type TokenType
    Lexeme string
    Literal Object
    Line int
}

// Returns a new "object"(read pointer) of type Token(read *Token)
func NewToken(tType TokenType, lexeme string, 
              literal Object, line int) Token {
    ret := Token {
        tType,
        lexeme,
        literal,
        line,
    }

    return ret
}
