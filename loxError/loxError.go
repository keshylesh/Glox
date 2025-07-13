package loxError

import (
    . "glox/util"
    . "glox/token"
    "os"
    "fmt"
)

type RuntimeError struct {
    Token Token
    Msg string
}

func (e *RuntimeError) Error() string {
    return fmt.Sprintf("%v - %v", e.Token, e.Msg)
}

func ErrorRuntime(error RuntimeError) {
    fmt.Fprintf(os.Stderr, "%v\n[line %v]\n", error.Msg, error.Token.Line)
    HadRuntimeError = true
}

type ReturnError struct {
    Value Object
}

func (e *ReturnError) Error() string {
    return fmt.Sprintf("%v", e.Value)
}
