package token

import (
    . "glox/util"
)

type Token struct {
    TType TokenType
    Lexeme string
    Literal Object
    Line int
}

func NewToken(tType TokenType, lexeme string, 
              literal Object, line int) *Token {
    ret := Token {
        tType,
        lexeme,
        literal,
        line,
    }

    return &ret
}
