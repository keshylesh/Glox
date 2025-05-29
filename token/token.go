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

func TokenError(token Token, msg string) {
    if token.Type == EOF {
        Report(token.Line, " at end", msg)
    } else {
        Report(token.Line, " at '" + token.Lexeme + "'", msg)
    }
}
