package ast

import (
	. "glox/token"
	. "glox/util"
)

type Visitor interface {
	VisitGrouping(obj Grouping) (Object, error)
	VisitLiteral(obj Literal) (Object, error)
	VisitUnary(obj Unary) (Object, error)
	VisitBinary(obj Binary) (Object, error)
}

type Expr interface{
	Accept(v Visitor) Object
}

type Literal struct {
	Value Object
}

func NewLiteral(Value Object) Literal {
	return Literal{Value,}
}

func (obj Literal) Accept(v Visitor) Object {
	ret, _ := v.VisitLiteral(obj)
	return ret
}

type Unary struct {
	Operator Token
	Right Expr
}

func NewUnary(Operator Token, Right Expr) Unary {
	return Unary{Operator, Right,}
}

func (obj Unary) Accept(v Visitor) Object {
	ret, _ := v.VisitUnary(obj)
	return ret
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

func NewBinary(Left Expr, Operator Token, Right Expr) Binary {
	return Binary{Left, Operator, Right,}
}

func (obj Binary) Accept(v Visitor) Object {
	ret, _ := v.VisitBinary(obj)
	return ret
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) Grouping {
	return Grouping{Expression,}
}

func (obj Grouping) Accept(v Visitor) Object {
	ret, _ := v.VisitGrouping(obj)
	return ret
}

