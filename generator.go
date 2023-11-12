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
		output += g.GenStmt(stmt) + "\n"
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
		output += "\tmov rax, 0\n" //set a default starting value
		output += g.push("rax")

	case NodeStmtVarAssign:
		variableName := stmt.ident.value.MustGetValue()
		variable, ok := g.variables[variableName]
		if !ok {
			panic(fmt.Errorf("variables must be declared before assignment. '%s' is undefined", variableName))
		}

		output += g.GenExpr(stmt.expr)
		output += g.pop("rax")
		output += fmt.Sprintf("\tmov QWORD [rsp + %v], rax\n", (g.stackSize-variable.stackLoc-1)*8)

	default:
		panic(fmt.Errorf("generator error: don't know how to generate statement: %T", rawStmt.variant))
	}
	return output
}

func (g *Generator) GenExpr(rawExpr NodeExpr) string {
	output := ""

	switch expr := rawExpr.variant.(type) {
	case NodeTerm:
		output += g.GenTerm(expr)
	case NodeBinExpr:
		output += g.GenBinExpr(expr)
	default:
		panic(fmt.Errorf("generator error: don't know how to generate expression: %T", rawExpr.variant))
	}
	return output
}

func (g *Generator) GenBinExpr(rawBinExpr NodeBinExpr) string {
	output := ""
	switch binExpr := rawBinExpr.variant.(type) {
	case NodeBinExprAdd:
		output += g.GenExpr(binExpr.left)
		output += g.GenExpr(binExpr.right)
		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tadd rax, rbx\n"
		output += g.push("rax")
	case NodeBinExprSubtract:
		output += g.GenExpr(binExpr.left)
		output += g.GenExpr(binExpr.right)
		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tsub rax, rbx\n"
		output += g.push("rax")
	case NodeBinExprMultiply:
		output += g.GenExpr(binExpr.left)
		output += g.GenExpr(binExpr.right)
		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tmul rbx\n"
		output += g.push("rax")
	case NodeBinExprDivide:
		output += g.GenExpr(binExpr.left)
		output += g.GenExpr(binExpr.right)
		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tdiv rbx\n"
		output += g.push("rax")
	default:
		panic(fmt.Errorf("generator error: don't know how to generate binary expression: %T", rawBinExpr.variant))
	}
	return output
}

func (g *Generator) GenTerm(rawTerm NodeTerm) string {
	output := ""

	switch term := rawTerm.variant.(type) {
	case NodeTermIntLiteral:
		output += "\tmov rax, " + term.intLiteral.value.MustGetValue() + "\n"
		output += g.push("rax")

	case NodeTermIdentifier:
		variableName := term.identifier.value.MustGetValue()

		variable, ok := g.variables[variableName]
		if !ok {
			panic(fmt.Errorf("unknown identifier: %v", variableName))
		}

		output += g.push(fmt.Sprintf("QWORD [rsp + %v]", (g.stackSize-variable.stackLoc-1)*8))

	default:
		panic(fmt.Errorf("generator error: don't know how to generate term: %T", rawTerm.variant))
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
