package ast

import (
	. "glox/token"
	. "glox/util"
)

type visitor interface {
	VisitBinary(obj Binary) Object
	VisitGrouping(obj Grouping) Object
	VisitLiteral(obj Literal) Object
	VisitUnary(obj Unary) Object
}

type Expr interface{
	Accept(v visitor) Object
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func NewBinary(left Expr, operator Token, right Expr) Binary {
	return Binary{left, operator, right,}
}

func (obj Binary) Accept(v visitor) Object {
	return v.VisitBinary(obj)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) Grouping {
	return Grouping{expression,}
}

func (obj Grouping) Accept(v visitor) Object {
	return v.VisitGrouping(obj)
}

type Literal struct {
	value Object
}

func NewLiteral(value Object) Literal {
	return Literal{value,}
}

func (obj Literal) Accept(v visitor) Object {
	return v.VisitLiteral(obj)
}

type Unary struct {
	operator Token
	right Expr
}

func NewUnary(operator Token, right Expr) Unary {
	return Unary{operator, right,}
}

func (obj Unary) Accept(v visitor) Object {
	return v.VisitUnary(obj)
}

