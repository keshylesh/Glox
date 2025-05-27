package ast

import (
    "fmt"
    . "glox/util"
)

type AstPrinter struct {
    v visitor
}

// Function to pretty print a tree
func (a AstPrinter) Print(expr Expr) string {
    ret := expr.Accept(a)
    return ret.(string)
}

// Defintions for visitor functions

func (a AstPrinter) VisitBinary(expr Binary) Object {
    return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (a AstPrinter) VisitGrouping(expr Grouping) Object {
    return a.parenthesize("group", expr.expression)
}

func (a AstPrinter) VisitLiteral(expr Literal) Object {
    if expr.value == nil {
        return "nil"
    }
    ret := fmt.Sprintf("%v", expr.value)
    return ret
}

func (a AstPrinter) VisitUnary(expr Unary) Object {
    return a.parenthesize(expr.operator.Lexeme, expr.right)
}


// Function to enclose parameters in brackets for ordering
func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
    ret := "(" + name
    for _, expr := range exprs {
        ret += " "
        temp := fmt.Sprintf("%v", expr.Accept(a))
        ret += temp
    }
    ret += ")"

    return ret
}
