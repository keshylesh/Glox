package interpreter

import (
    . "glox/ast"
    . "glox/token"
    . "glox/util"
    "reflect"
    "fmt"
    "errors"
    "os"
)

type RuntimeError struct {
    token Token
    msg string
}

func (e *RuntimeError) Error() string {
    return fmt.Sprintf("%v - %v", e.token, e.msg)
}

func ErrorRuntime(error RuntimeError) {
    fmt.Fprintf(os.Stderr, "%v\n[line %v]\n", error.msg, error.token.Line)
    HadRuntimeError = true
}


type Interpreter struct {
    v Visitor
}

func (i Interpreter) Interpret(expression Expr) {
    val, err := i.evaluate(expression)
    var re *RuntimeError
    if errors.As(err, &re) {
        ErrorRuntime(*re)
        return
    }
    fmt.Println(stringify(val))
}

func (i Interpreter) VisitLiteral(expr Literal) (Object, error) {
    return expr.Value, nil
}

func (i Interpreter) VisitGrouping(expr Grouping) (Object, error) {
    return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitUnary(expr Unary) (Object, error) {
    right, err := i.evaluate(expr.Right)
    if err != nil {
        return nil, err
    }

    switch expr.Operator.Type {
    case BANG:
        return !isTruthy(right), nil
    case MINUS:
        err = verifyType("float64", "number", expr.Operator, right)
        if err != nil {
            return nil, err
        }
        return -(right.(float64)), nil
    }

    return nil, nil
}

func (i Interpreter) VisitBinary(expr Binary) (Object, error) {
    left, err := i.evaluate(expr.Left)
    if err != nil {
        return nil, err
    }
    right, err := i.evaluate(expr.Right)
    if err != nil {
        return nil, err
    }

    switch expr.Operator.Type {
    case GREAT:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) > right.(float64), nil
        
    case GREAT_EQUAL:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) >= right.(float64), nil

    case LESS:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) < right.(float64), nil

    case LESS_EQUAL:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) <= right.(float64), nil

    case BANG_EQUAL:
        return !isEqual(left, right), nil

    case EQUAL_EQUAL:
        return isEqual(left, right), nil

    case MINUS:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) - right.(float64), nil

    case PLUS:
        if typeOf(left) == "float64" && typeOf(right) == "float64" {
            return left.(float64) + right.(float64), nil
        }

        if typeOf(left) == "string" || typeOf(right) == "string" {
            if typeOf(left) == "float64" {
                left = fmt.Sprintf("%v", left.(float64))
            }
            if typeOf(right) == "float64" {
                right = fmt.Sprintf("%v", right.(float64))
            }
            return left.(string) + right.(string), nil
        }

        return nil, &RuntimeError{expr.Operator, "Operand(s) must be two numbers or two strings"}

    case SLASH:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        if right.(float64) == 0 {
            return nil, &RuntimeError{expr.Operator, "Cannot divide by zero"}
        }
        return left.(float64) / right.(float64), nil

    case STAR:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) * right.(float64), nil
    }

    return nil, nil
}

func (i Interpreter) evaluate(expr Expr) (Object, error) {
    return expr.Accept(i)
}

func isTruthy(obj Object) bool {
    if obj == nil {
        return false
    }
    switch obj.(type) {
    case bool:
        return obj.(bool)
    default:
        return true
    }
}

func isEqual(x, y Object) bool {
    if x == nil && y == nil {
        return true
    } 
    if x == nil {
        return false
    }
    return reflect.DeepEqual(x, y)
}

func stringify(obj Object) string {
    if obj == nil {
        return "nil"
    } 

    return fmt.Sprintf("%v", obj) 
}

func typeOf(obj Object) string {
    return reflect.TypeOf(obj).String()
}

func verifyType(match, eMsg string, operator Token, operands ...Object) error {
    for _, operand := range operands {
        if (typeOf(operand) != match) {
            return &RuntimeError{operator, "Operand(s) must be a " + eMsg}
        }
    }
    
    return nil
}
