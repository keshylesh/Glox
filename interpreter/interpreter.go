package interpreter

import (
    . "glox/ast"
    . "glox/token"
    . "glox/util"
    . "glox/environment"
    . "glox/loxError"
    "reflect"
    "fmt"
    "errors"
)

type Interpreter struct {
    ev ExprVisitor
    sv StmtVisitor
    env *Environment
    globals *Environment
}

// Interpreter "constructor"
func NewInterpreter() Interpreter {
    global := NewEnvironment()
    var clock Clock 
    global.Define("clock", clock)

    return Interpreter{env: global, globals: global}
}

// function to interpret a series of statements
func (i Interpreter) Interpret(statements []Stmt) {
    for _, statement := range statements {
        err := i.execute(statement)
        var re *RuntimeError
        if errors.As(err, &re) {
            ErrorRuntime(*re)
            return
        }
    }
}

// VISTITOR FUNCTIONS

func (i Interpreter) VisitLiteral(expr Literal) (Object, error) {
    return expr.Value, nil
}

func (i Interpreter) VisitLogical(expr Logical) (Object, error) {
    left, err := i.evaluate(expr.Left)
    if err != nil { return nil, err }

    if expr.Operator.Type == AND && !isTruthy(left) {
        return left, nil 
    }
    if expr.Operator.Type == OR && isTruthy(left) {
        return left, nil 
    }

    return i.evaluate(expr.Right)
}

func (i Interpreter) VisitGrouping(expr Grouping) (Object, error) {
    return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitUnary(expr Unary) (Object, error) {
    right, err := i.evaluate(expr.Right)
    if err != nil { return nil, err }

    switch expr.Operator.Type {
    case BANG:
        return !isTruthy(right), nil
    case MINUS:
        err = verifyType("float64", "number", expr.Operator, right)
        if err != nil { return nil, err }

        return -(right.(float64)), nil
    }

    return nil, nil
}

func (i Interpreter) VisitVariable(expr Variable) (Object, error) {
    return i.env.Get(expr.Name)
}

func (i Interpreter) VisitBinary(expr Binary) (Object, error) {
    left, err := i.evaluate(expr.Left)
    if err != nil { return nil, err }

    right, err := i.evaluate(expr.Right)
    if err != nil { return nil, err }

    switch expr.Operator.Type {
    case GREAT:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil { return nil, err }

        return left.(float64) > right.(float64), nil
        
    case GREAT_EQUAL:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil { return nil, err }

        return left.(float64) >= right.(float64), nil

    case LESS:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil {
            return nil, err
        }
        return left.(float64) < right.(float64), nil

    case LESS_EQUAL:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil { return nil, err }

        return left.(float64) <= right.(float64), nil

    case BANG_EQUAL:
        return !isEqual(left, right), nil

    case EQUAL_EQUAL:
        return isEqual(left, right), nil

    case MINUS:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil { return nil, err }

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
        if err != nil { return nil, err }

        if right.(float64) == 0 {
            return nil, &RuntimeError{expr.Operator, "Cannot divide by zero"}
        }

        return left.(float64) / right.(float64), nil

    case STAR:
        err = verifyType("float64", "number", expr.Operator, right, left)
        if err != nil { return nil, err }

        return left.(float64) * right.(float64), nil
    }

    return nil, nil
}

func (i Interpreter) VisitCall(expr Call) (Object, error) {
    callee, err := i.evaluate(expr.Callee)
    if err != nil { return nil, err }

    args := make([]Object, 0)
    for _, arg := range expr.Arguments {
        val, err := i.evaluate(arg)
        if err != nil { return nil, err }
        args = append(args, val)
    }

    function, ok := callee.(Callable)
    if !ok {
        return nil, &RuntimeError{expr.Paren, "Can only call functions and classes"}
    }
    if len(args) != function.Arity() {
        errMsg := fmt.Sprintf("Expected %v but got %v", function.Arity(), len(args))
        return nil, &RuntimeError{expr.Paren, errMsg}
    }

    return function.Call(i, args)
}

func (i Interpreter) evaluate(expr Expr) (Object, error) {
    return expr.Accept(i)
}

func (i Interpreter) execute(stmt Stmt) error {
    _, err := stmt.Accept(i)
    return err
}

func (i Interpreter) executeBlock(statements []Stmt, env *Environment) error {
    previous := i.env
    i.env = env

    for _, stmt := range statements {
        err := i.execute(stmt)
        if err != nil {
            i.env = previous
            return err
        }
    }

    i.env = previous
    return nil
}

func (i Interpreter) VisitBlock(stmt Block) (Object, error) {
    err := i.executeBlock(stmt.Statements, NewEnvironment(i.env))
    if err != nil { return nil, err }
    return nil, nil
}

func (i Interpreter) VisitStmtExpression(stmt StmtExpression) (Object, error) {
    _, err := i.evaluate(stmt.Expression)
    return nil, err
}

func (i Interpreter) VisitFunction(stmt Function) (Object, error) {
    function := NewLoxFunction(stmt, i.env)
    i.env.Define(stmt.Name.Lexeme, function)
    return nil, nil
}

func (i Interpreter) VisitIf(stmt If) (Object, error) {
    cond, err := i.evaluate(stmt.Condition)
    if err != nil { return nil, err }

    if isTruthy(cond) {
        err = i.execute(stmt.ThenBranch)
        if err != nil { return nil, err }
    } else if stmt.ElseBranch != nil {
        err = i.execute(stmt.ElseBranch)
        if err != nil { return nil, err }
    }

    return nil, nil
}

func (i Interpreter) VisitPrint(stmt Print) (Object, error) {
    val, err := i.evaluate(stmt.Expression)
    if err == nil {
        fmt.Println(stringify(val))
    }

    return nil, err
}

func (i Interpreter) VisitReturn(stmt Return) (Object, error) {
    var val Object = nil
    if stmt.Value != nil {
        value, err := i.evaluate(stmt.Value) 
        val = value
        if err != nil { 
            return nil, err
        }
    }

    return nil, &ReturnError{ val }
}

func (i Interpreter) VisitVar(stmt Var) (Object, error) {
    var val Object = nil
    if stmt.Initializer != nil {
        value, err := i.evaluate(stmt.Initializer)
        val = value
        if err != nil {
            return nil, err
        }
    }

    i.env.Define(stmt.Name.Lexeme, val)
    return nil, nil
}

func (i Interpreter) VisitWhile(stmt While) (Object, error) {
    cond, err := i.evaluate(stmt.Condition)
    if err != nil { return nil, err }

    for isTruthy(cond) {
        err = i.execute(stmt.Body)
        if err != nil { return nil, err }

        cond, err = i.evaluate(stmt.Condition)
        if err != nil { return nil, err }
    }
    
    return nil, nil
}

func (i Interpreter) VisitAssign(expr Assign) (Object, error) {
    value, err := i.evaluate(expr.Value)
    if err != nil { return nil, err }

    i.env.Assign(expr.Name, value)
    return value, nil
}


// function to return whether an Object is a truth-like value
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

// function to return whether two objects are equal
func isEqual(x, y Object) bool {
    if x == nil && y == nil {
        return true
    } 

    if x == nil {
        return false
    }

    return reflect.DeepEqual(x, y)
}

// function to turn an Object to a string representation
func stringify(obj Object) string {
    if obj == nil {
        return "nil"
    }

    // HACK: Go has no innate "ToString" manually checking if object is callable
    if function, ok := obj.(Callable); ok {
        return function.ToString()
    }

    return fmt.Sprintf("%v", obj) 
}

// function to return the type of an Object
func typeOf(obj Object) string {
    return reflect.TypeOf(obj).String()
}

// function to verify that the type of all objects are the passed in type
func verifyType(match, eMsg string, operator Token, operands ...Object) error {
    for _, operand := range operands {
        if (typeOf(operand) != match) {
            return &RuntimeError{operator, "Operand(s) must be a " + eMsg}
        }
    }
    
    return nil
}
