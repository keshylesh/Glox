package ast

import (
	. "glox/util"
)

type StmtVisitor interface {
	VisitStmtExpression(obj StmtExpression) (Object, error)
	VisitPrint(obj Print) (Object, error)
}

type Stmt interface{
	Accept(v StmtVisitor) (Object, error)
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

type StmtExpression struct {
	Expression Expr
}

func NewStmtExpression(Expression Expr) StmtExpression {
	return StmtExpression{Expression,}
}

func (obj StmtExpression) Accept(v StmtVisitor) (Object, error) {
	return v.VisitStmtExpression(obj)
}

