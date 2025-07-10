package interpreter

import (
    . "glox/util"
    "time"
)

type Clock struct {}

func (c Clock) Arity() int {
    return 0
}

func (c Clock) Call(i Interpreter, args []Object) (Object, error) {
    return float64(time.Now().UnixMilli()) / 1000.0, nil 
}

func (c Clock) ToString() string {
    return "<native fn>"
}
