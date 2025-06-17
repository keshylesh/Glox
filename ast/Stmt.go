package ast

import (
	. "glox/token"
	. "glox/util"
)

type StmtVisitor interface {
	VisitPrint(obj Print) (Object, error)
	VisitVar(obj Var) (Object, error)
	VisitStmtExpression(obj StmtExpression) (Object, error)
}

type Stmt interface{
	Accept(v StmtVisitor) (Object, error)
}

type StmtExpression struct {
	Expression Expr
}

func NewStmtExpression(Expression Expr) StmtExpression {
	return StmtExpression{Expression,}
}

func (obj StmtExpression) Accept(v StmtVisitor) (Object, error) {
	return v.VisitStmtExpression(obj)
}

type Print struct {
	Expression Expr
}

func NewPrint(Expression Expr) Print {
	return Print{Expression,}
}

func (obj Print) Accept(v StmtVisitor) (Object, error) {
	return v.VisitPrint(obj)
}

type Var struct {
	Name Token
	Initializer Expr
}

func NewVar(Name Token, Initializer Expr) Var {
	return Var{Name, Initializer,}
}

func (obj Var) Accept(v StmtVisitor) (Object, error) {
	return v.VisitVar(obj)
}

