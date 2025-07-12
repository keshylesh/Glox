package interpreter

import (
	. "glox/ast"
	. "glox/util"
	. "glox/environment"
)

type LoxFunction struct {
    declaration Function
}

func NewLoxFunction(decl Function) *LoxFunction {
    return &LoxFunction{ declaration: decl }
}

func (f LoxFunction) Call(i Interpreter, args []Object) (Object, error) {
    env := NewEnvironment(i.globals)
    for k := 0; k < len(f.declaration.Params); k++ {
        env.Define(f.declaration.Params[k].Lexeme, args[k])
    }

    err := i.executeBlock(f.declaration.Body, env)
    if err != nil { return nil, err }

    return nil, nil
}

func (f LoxFunction) Arity() int {
    return len(f.declaration.Params)
}

func (f LoxFunction) ToString() string {
    return "<fn " + f.declaration.Name.Lexeme + ">"
}
