package ast

import (
	. "glox/token"
	. "glox/util"
)

type ExprVisitor interface {
	VisitBinary(obj Binary) (Object, error)
	VisitGrouping(obj Grouping) (Object, error)
	VisitLiteral(obj Literal) (Object, error)
	VisitUnary(obj Unary) (Object, error)
}

type Expr interface{
	Accept(v ExprVisitor) (Object, error)
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

func NewBinary(Left Expr, Operator Token, Right Expr) Binary {
	return Binary{Left, Operator, Right,}
}

func (obj Binary) Accept(v ExprVisitor) (Object, error) {
	return v.VisitBinary(obj)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) Grouping {
	return Grouping{Expression,}
}

func (obj Grouping) Accept(v ExprVisitor) (Object, error) {
	return v.VisitGrouping(obj)
}

type Literal struct {
	Value Object
}

func NewLiteral(Value Object) Literal {
	return Literal{Value,}
}

func (obj Literal) Accept(v ExprVisitor) (Object, error) {
	return v.VisitLiteral(obj)
}

type Unary struct {
	Operator Token
	Right Expr
}

func NewUnary(Operator Token, Right Expr) Unary {
	return Unary{Operator, Right,}
}

func (obj Unary) Accept(v ExprVisitor) (Object, error) {
	return v.VisitUnary(obj)
}

