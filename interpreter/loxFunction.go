package interpreter

import (
	. "glox/ast"
	. "glox/util"
	. "glox/environment"
	. "glox/loxError"
	"errors"
)

type LoxFunction struct {
    declaration Function
    closure *Environment
}

func NewLoxFunction(decl Function, closure *Environment) *LoxFunction {
    return &LoxFunction{ declaration: decl, closure: closure }
}

func (f LoxFunction) Call(i Interpreter, args []Object) (Object, error) {
    env := NewEnvironment(f.closure)
    for k := 0; k < len(f.declaration.Params); k++ {
        env.Define(f.declaration.Params[k].Lexeme, args[k])
    }

    err := i.executeBlock(f.declaration.Body, env)
    var re *ReturnError
    if errors.As(err, &re) {
        return re.Value, nil
    } else if err != nil { 
        return nil, err
    }

    return nil, nil
}

func (f LoxFunction) Arity() int {
    return len(f.declaration.Params)
}

func (f LoxFunction) ToString() string {
    return "<fn " + f.declaration.Name.Lexeme + ">"
}
