package ast

import (
    "fmt"
    . "glox/util"
)

type AstPrinter struct {
    v ExprVisitor
}

// Function to pretty print a tree
func (a AstPrinter) Print(expr Expr) string {
    ret, _ := expr.Accept(a)
    return ret.(string)
}

// Defintions for visitor functions

func (a AstPrinter) VisitBinary(expr Binary) (Object, error) {
    return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a AstPrinter) VisitGrouping(expr Grouping) (Object, error) {
    return a.parenthesize("group", expr.Expression), nil
}

func (a AstPrinter) VisitLiteral(expr Literal) (Object, error) {
    if expr.Value == nil {
        return "nil", nil
    }
    ret := fmt.Sprintf("%v", expr.Value)
    return ret, nil
}

func (a AstPrinter) VisitUnary(expr Unary) (Object, error) {
    return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

// Function to enclose parameters in brackets for ordering
func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
    ret := "(" + name
    for _, expr := range exprs {
        ret += " "
        val, _ := expr.Accept(a)
        temp := fmt.Sprintf("%v", val)
        ret += temp
    }
    ret += ")"

    return ret
}
