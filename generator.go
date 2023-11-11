package main

import (
	"fmt"
)

type Generator struct {
	program   NodeProg
	stackSize uint
	variables map[string]Variable
}

func NewGenerator(prog NodeProg) Generator {
	return Generator{
		program:   prog,
		stackSize: 0,
		variables: map[string]Variable{},
	}
}

func (g *Generator) GenProg() string {
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

func (g *Generator) GenStmt(rawStmt NodeStmt) string {
	output := ""

	switch stmt := rawStmt.variant.(type) {
	case NodeStmtExit:
		output += g.GenExpr(stmt.expr)
		output += "\tmov rax, 60\n"
		output += g.pop("rdi")
		output += "\tsyscall\n"
	case NodeStmtVarDeclare:
		variableName := stmt.ident.value.MustGetValue()

		if _, ok := g.variables[variableName]; ok {
			panic(fmt.Errorf("identifier already used: %v", variableName))
		}
		g.variables[variableName] = Variable{stackLoc: g.stackSize}
		output += "\tmov rbx, 0\n" //set a default starting value
		output += g.push("rbx")

	case NodeStmtVarAssign:
		// variableName := stmt.ident.value.MustGetValue()

		// output += g.GenExpr(stmt.expr)

	default:
		panic(fmt.Errorf("generator error: don't know how to generate statement: %T", rawStmt.variant))
	}
	return output
}

func (g *Generator) GenExpr(rawExpr NodeExpr) string {
	output := ""

	switch expr := rawExpr.variant.(type) {
	case NodeExprIntLiteral:
		output += "\tmov rax, " + expr.intLiteral.value.MustGetValue() + "\n"
		output += g.push("rax")
	case NodeExprIdentifier:
		variableName := expr.identifier.value.MustGetValue()

		variable, ok := g.variables[variableName]
		if !ok {
			panic(fmt.Errorf("unknown identifier: %v", variableName))
		}

		output += g.push(fmt.Sprintf("QWORD [rsp + %v]", (g.stackSize-variable.stackLoc-1)*8))
	default:
		panic(fmt.Errorf("generator error: don't know how to generate expression: %T", rawExpr.variant))
	}
	return output
}

func (g *Generator) push(reg string) string {
	g.stackSize++
	return "\tpush " + reg + "\n"
}
func (g *Generator) pop(reg string) string {
	g.stackSize--
	return "\tpop " + reg + "\n"
}

type Variable struct {
	stackLoc uint
}
