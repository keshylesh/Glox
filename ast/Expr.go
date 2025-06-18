package ast

import (
	. "glox/token"
	. "glox/util"
)

type ExprVisitor interface {
	VisitLogical(obj Logical) (Object, error)
	VisitUnary(obj Unary) (Object, error)
	VisitVariable(obj Variable) (Object, error)
	VisitAssign(obj Assign) (Object, error)
	VisitBinary(obj Binary) (Object, error)
	VisitGrouping(obj Grouping) (Object, error)
	VisitLiteral(obj Literal) (Object, error)
}

type Expr interface{
	Accept(v ExprVisitor) (Object, error)
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

type Logical struct {
	Left Expr
	Operator Token
	Right Expr
}

func NewLogical(Left Expr, Operator Token, Right Expr) Logical {
	return Logical{Left, Operator, Right,}
}

func (obj Logical) Accept(v ExprVisitor) (Object, error) {
	return v.VisitLogical(obj)
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

type Variable struct {
	Name Token
}

func NewVariable(Name Token) Variable {
	return Variable{Name,}
}

func (obj Variable) Accept(v ExprVisitor) (Object, error) {
	return v.VisitVariable(obj)
}

type Assign struct {
	Name Token
	Value Expr
}

func NewAssign(Name Token, Value Expr) Assign {
	return Assign{Name, Value,}
}

func (obj Assign) Accept(v ExprVisitor) (Object, error) {
	return v.VisitAssign(obj)
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

