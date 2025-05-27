package util

import (
    "fmt"
    "os"
)

// Global var used to stop line/program from being run
var HadError bool = false

// Panics when encountering an error
func Check(e error) {
    if e != nil {
        panic(e)
    }
}

// Report an error with the given line number and message
func Error(line int, msg string) {
    Report(line, "", msg)
}

// Print out line error to stderr
// no panic() as multiple errors in multiple lines are valid
func Report(line int, where string, msg string) {
    fmt.Fprintf(os.Stderr, "[line %v] Error %v: %v", line, where, msg)

    // Ensure program doesn't run (for main.runFile())
    // Ensure line doesn't run (for main.runPrompt())
    HadError = true
}

// Util functions to check if characters are alphabets/digits

func IsAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func IsAlphaNumeric(c byte) bool {
    return IsAlpha(c) || IsDigit(c)
}

func IsDigit(c byte) bool {
    return c >= '0' && c <= '9'
}
