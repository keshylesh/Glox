package scanner

import (
    . "glox/util"
    . "glox/token"
    "strconv"
)

type Scanner struct {
    src string
    tokens []Token
    start int
    current int
    line int
}

func NewScanner(src string) *Scanner {
    ret := Scanner {
        src: src,
        line: 1,
    }

    return &ret
}

func (s *Scanner) ScanTokens() []Token {
    for !s.isAtEnd() {
        s.start = s.current
        s.scanToken()
    }

    s.tokens = append(s.tokens, *NewToken(EOF, "", nil, s.line))
    return s.tokens
}

func (s *Scanner) scanToken() {
    c := s.advance()
    switch c {
    case '(':
        s.addToken(LEFT_PAREN, nil)
    case ')':
        s.addToken(RIGHT_PAREN, nil)
    case '{':
        s.addToken(LEFT_BRACE, nil)
    case '}':
        s.addToken(RIGHT_BRACE, nil)
    case ',':
        s.addToken(COMMA, nil)
    case '.':
        s.addToken(DOT, nil)
    case '-':
        s.addToken(MINUS, nil)
    case '+':
        s.addToken(PLUS, nil)
    case ';':
        s.addToken(SEMICOLON, nil)
    case '*':
        s.addToken(STAR, nil)
    case '!':
        if s.match('=') {
            s.addToken(BANG_EQUAL, nil)
        } else {
            s.addToken(BANG, nil)
        }
    case '=':
        if s.match('=') {
            s.addToken(EQUAL_EQUAL, nil)
        } else {
            s.addToken(EQUAL, nil)
        }
    case '<':
        if s.match('=') {
            s.addToken(LESS_EQUAL, nil)
        } else {
            s.addToken(LESS, nil)
        }
    case '>':
        if s.match('=') {
            s.addToken(GREAT_EQUAL, nil)
        } else {
            s.addToken(GREAT, nil)
        }
    case '/':
        if s.match('/') {
            for s.peek() != '\n' && !s.isAtEnd() {
                s.advance()
            }
        } else {
            s.addToken(SLASH, nil)
        }
    case ' ':
        fallthrough
    case '\r':
        fallthrough
    case '\t':
    case '\n':
        s.line++
    case '"':
        s.string()
    default:
        if IsDigit(c) {
            s.number()
        } else if IsAlpha(c) {
            s.identifier()
        } else {
            Error(s.line, "Unexpected character.")
        }
    }
}

func (s *Scanner) identifier() {
    for IsAlphaNumeric(s.peek()) {
        s.advance()
    }

    text := s.src[s.start:s.current]
    tType := Keywords[text]
    if tType == NO_TYPE {
        tType = IDENTIFIER
    }
    
    s.addToken(tType, nil)
}

func (s *Scanner) string() {
    for s.peek() != '"' && !s.isAtEnd() {
        if s.peek() == '\n' {
            s.line++
        }
        s.advance()
    }

    if s.isAtEnd() {
        Error(s.line, "Unterminated string.")
    }

    // Eat closing quote
    s.advance()

    // Trim quotes
    value := s.src[s.start + 1:s.current - 1]
    s.addToken(STRING, value)
}

func (s *Scanner) number() {
    for IsDigit(s.peek()) {
        s.advance()
    }

    // Look for a fraction
    if s.peek() == '.' && IsDigit(s.peekNext()) {
        s.advance()

        for IsDigit(s.peek()) {
            s.advance()
        }
    }

    f, e := strconv.ParseFloat(s.src[s.start : s.current], 64)
    if e != nil {
        panic(e)
    }
    s.addToken(NUMBER, f)
}

func (s *Scanner) match(expected byte) bool {
    if s.isAtEnd() || s.src[s.current] != expected {
        return false
    }

    s.current++
    return true
} 

func (s *Scanner) peek() byte {
    if s.isAtEnd() {
        return 0
    }

    return s.src[s.current]
}

func (s *Scanner) peekNext() byte {
    if s.current + 1 >= len(s.src) {
        return 0
    }

    return s.src[s.current + 1]
}

func (s *Scanner) isAtEnd() bool {
    return s.current >= len(s.src)
}

func (s *Scanner) advance() byte {
    ret := s.src[s.current]
    s.current++
    return ret
}

func (s *Scanner) addToken(tType TokenType, literal Object) {
    text := s.src[s.start: s.current]
    s.tokens = append(s.tokens, *NewToken(tType, text, literal, s.line))
}
