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

// Return a new "object"(read pointer) of type Scanner(read *Scanner) 
func NewScanner(src string) *Scanner {
    ret := Scanner {
        src: src,
        line: 1,
    }

    return &ret
}

// Public method of Scanner that scans the src string for tokens
// and returns the token slice
func (s *Scanner) ScanTokens() []Token {
    for !s.isAtEnd() {
        s.start = s.current
        s.scanToken()
    }

    s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
    return s.tokens
}

// Method of Scanner to scan a single token
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
        // Discard whitespace
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

// Function to finish scanning an identifier. It checks if the text maps to 
// an existing TokenType else is NO_TYPE and assigns it as IDENTIFIER
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

// Function to finish scanning a string 
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

// Function to finish scanning a number
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

    // Convert string -> double
    f, err := strconv.ParseFloat(s.src[s.start : s.current], 64)
    Check(err)
    s.addToken(NUMBER, f)
}

// Function to check if the next byte matches the expected byte
func (s *Scanner) match(expected byte) bool {
    if s.isAtEnd() || s.src[s.current] != expected {
        return false
    }

    s.current++
    return true
} 

// Function to peek at the next byte without advancing
func (s *Scanner) peek() byte {
    if s.isAtEnd() {
        return 0
    }

    return s.src[s.current]
}

// Function to peek at the next to next byte without advancing
func (s *Scanner) peekNext() byte {
    if s.current + 1 >= len(s.src) {
        return 0
    }

    return s.src[s.current + 1]
}

// Function to check if source string is at the end
func (s *Scanner) isAtEnd() bool {
    return s.current >= len(s.src)
}

// Function to return the current byte and move the current index forward
func (s *Scanner) advance() byte {
    ret := s.src[s.current]
    s.current++
    return ret
}

// Function to add a token to the tokens slice 
func (s *Scanner) addToken(tType TokenType, literal Object) {
    text := s.src[s.start: s.current]
    s.tokens = append(s.tokens, NewToken(tType, text, literal, s.line))
}
