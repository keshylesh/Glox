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
    enclosing *Environment
    values map[string]Object
}

// Environment "constructor"
func NewEnvironment(params ...*Environment) *Environment {
    //HACK: mimicing default values/overloading with variadic function
    var enclosing *Environment = nil
    if len(params) > 0 {
        enclosing = params[0]
    }
    return &Environment{enclosing: enclosing, values: make(map[string]Object)}
}

// function to define a new variable with a value
func (e *Environment) Define(name string, value Object) {
    e.values[name] = value
}

// function to retrieve the value associated with a given variable name
// recursively check the enclosing scope for the variable if not found
func (e *Environment) Get(name Token) (Object, error) {
    if val, ok := e.values[name.Lexeme]; ok {
        return val, nil
    }

    if e.enclosing != nil {
        return e.enclosing.Get(name)
    }

    return nil, &RuntimeError{name, "Undefined variable '" + name.Lexeme + "'"}
}

// function to assign a value to an existing variable
// recursively check the enclosing scope for the variable if not found
func (e *Environment) Assign(name Token, value Object) error {
    if _, ok := e.values[name.Lexeme]; ok {
        e.values[name.Lexeme] = value
        return nil
    }

    if e.enclosing != nil {
        return e.enclosing.Assign(name, value)
    }
    
    return &RuntimeError{name, "Undefined variable '" + name.Lexeme + "'"}
}

