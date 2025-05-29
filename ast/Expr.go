package ast

import (
	. "glox/token"
	. "glox/util"
)

type visitor interface {
	VisitUnary(obj Unary) Object
	VisitBinary(obj Binary) Object
	VisitGrouping(obj Grouping) Object
	VisitLiteral(obj Literal) Object
}

type Expr interface{
	Accept(v visitor) Object
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

func NewBinary(Left Expr, Operator Token, Right Expr) Binary {
	return Binary{Left, Operator, Right,}
}

func (obj Binary) Accept(v visitor) Object {
	return v.VisitBinary(obj)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) Grouping {
	return Grouping{Expression,}
}

func (obj Grouping) Accept(v visitor) Object {
	return v.VisitGrouping(obj)
}

type Literal struct {
	Value Object
}

func NewLiteral(Value Object) Literal {
	return Literal{Value,}
}

func (obj Literal) Accept(v visitor) Object {
	return v.VisitLiteral(obj)
}

type Unary struct {
	Operator Token
	Right Expr
}

func NewUnary(Operator Token, Right Expr) Unary {
	return Unary{Operator, Right,}
}

func (obj Unary) Accept(v visitor) Object {
	return v.VisitUnary(obj)
}

