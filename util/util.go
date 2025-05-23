package util

import (
    "fmt"
    "os"
)

var HadError bool = false

func Error(line int, msg string) {
    Report(line, "", msg)
}

func Report(line int, where string, msg string) {
    fmt.Fprintf(os.Stderr, "[line %v] Error %v: %v", line, where, msg)
    HadError = true
}

func IsAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func IsAlphaNumeric(c byte) bool {
    return IsAlpha(c) || IsDigit(c)
}

func IsDigit(c byte) bool {
    return c >= '0' && c <= '9'
}
