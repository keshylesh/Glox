package ast

import (
	. "glox/token"
	. "glox/util"
)

type StmtVisitor interface {
	VisitStmtExpression(obj StmtExpression) (Object, error)
	VisitFunction(obj Function) (Object, error)
	VisitIf(obj If) (Object, error)
	VisitPrint(obj Print) (Object, error)
	VisitVar(obj Var) (Object, error)
	VisitWhile(obj While) (Object, error)
	VisitBlock(obj Block) (Object, error)
}

type Stmt interface{
	Accept(v StmtVisitor) (Object, error)
}

type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(Condition Expr, ThenBranch Stmt, ElseBranch Stmt) If {
	return If{Condition, ThenBranch, ElseBranch,}
}

func (obj If) Accept(v StmtVisitor) (Object, error) {
	return v.VisitIf(obj)
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

type While struct {
	Condition Expr
	Body Stmt
}

func NewWhile(Condition Expr, Body Stmt) While {
	return While{Condition, Body,}
}

func (obj While) Accept(v StmtVisitor) (Object, error) {
	return v.VisitWhile(obj)
}

type Block struct {
	Statements []Stmt
}

func NewBlock(Statements []Stmt) Block {
	return Block{Statements,}
}

func (obj Block) Accept(v StmtVisitor) (Object, error) {
	return v.VisitBlock(obj)
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

type Function struct {
	Name Token
	Params []Token
	Body []Stmt
}

func NewFunction(Name Token, Params []Token, Body []Stmt) Function {
	return Function{Name, Params, Body,}
}

func (obj Function) Accept(v StmtVisitor) (Object, error) {
	return v.VisitFunction(obj)
}

