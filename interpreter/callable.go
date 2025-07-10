package interpreter

import (
    . "glox/util"
)

type Callable interface {
    Arity() int
    Call(interpreter Interpreter, args []Object) (Object, error)
}
