package environment

import (
    . "glox/util"
    . "glox/token"
    "fmt"
)

type RuntimeError struct {
    Token Token
    Msg string
}

func (e *RuntimeError) Error() string {
    return fmt.Sprintf("%v - %v", e.Token, e.Msg)
}

type Environment struct {
    values map[string]Object
}

func NewEnvironment() Environment {
    return Environment{make(map[string]Object)}
}

func (e Environment) Define(name string, value Object) {
    e.values[name] = value
}

func (e Environment) Get(name Token) (Object, error) {
    if val, ok := e.values[name.Lexeme]; ok {
        return val, nil
    }

    return nil, &RuntimeError{name, "Undefined variable '" + name.Lexeme + "'"}
}

func (e Environment) Assign(name Token, value Object) error {
    if _, ok := e.values[name.Lexeme]; ok {
        e.values[name.Lexeme] = value
        return nil
    }
    
    return &RuntimeError{name, "Undefined variable '" + name.Lexeme + "'"}
}

