package main

import (
	"fmt"
)

type Generator struct {
	program NodeProg
}

func NewGenerator(prog NodeProg) Generator {
	return Generator{
		program: prog,
	}
}

func (g Generator) GenProg() string {
	output := "global _start\n_start:\n"

	for _, stmt := range g.program.stmts {
		output += g.GenStmt(stmt)
	}

	//exit 0 at end of program if no explicit exit called
	output += "\tmov rax, 60\n"
	output += "\tmov rdi, 0\n"
	output += "\tsyscall\n"

	return output
}

func (g Generator) GenStmt(rawStmt NodeStmt) string {
	output := ""

	switch stmt := rawStmt.variant.(type) {
	case NodeStmtExit:
		output += g.GenExpr(stmt.expr)
		output += "\tmov rax, 60\n"
		output += "\tpop rdi\n"
		output += "\tsyscall\n"
	case NodeStmtVarAssign:
	case NodeStmtVarDeclare:
	default:
		panic(fmt.Errorf("generator error: don't know how to generate statement: %T", rawStmt.variant))
	}
	return output
}

func (g Generator) GenExpr(rawExpr NodeExpr) string {
	output := ""

	switch expr := rawExpr.variant.(type) {
	case NodeExprIntLiteral:
		output += "\tmov rax, " + expr.intLiteral.value.MustGetValue() + "\n"
		output += "\tpush rax\n"
	case NodeExprIdentifier:
		//TODO
	default:
		panic(fmt.Errorf("generator error: don't know how to generate expression: %T", rawExpr.variant))
	}
	return output
}
