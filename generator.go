package main

import (
	"fmt"
)

type Generator struct {
	program    NodeProg
	stackSize  uint
	variables  []Variable
	scopes     []int
	labelCount int
}

func NewGenerator(prog NodeProg) Generator {
	return Generator{
		program:    prog,
		stackSize:  0,
		variables:  []Variable{},
		scopes:     []int{},
		labelCount: 0,
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

		exists := false
		for _, v := range g.variables {
			if v.name == variableName {
				exists = true
			}
		}

		if exists {
			panic(fmt.Errorf("identifier already used: %v", variableName))
		}
		g.variables = append(g.variables, Variable{stackLoc: g.stackSize, name: variableName})
		output += "\tmov rax, 0\n" //set a default starting value
		output += g.push("rax")

	case NodeStmtVarAssign:
		variableName := stmt.ident.value.MustGetValue()
		var variable Variable
		exists := false
		for _, v := range g.variables {
			if v.name == variableName {
				variable = v
				exists = true
			}
		}
		if !exists {
			panic(fmt.Errorf("variables must be declared before assignment. '%s' is undefined", variableName))
		}

		output += g.GenExpr(stmt.expr)
		output += g.pop("rax")
		output += fmt.Sprintf("\tmov QWORD [rsp + %v], rax\n", (g.stackSize-variable.stackLoc-1)*8)

	case NodeScope:
		output += g.GenScope(stmt)
	case NodeStmtIf:
		output += g.GenExpr(stmt.expr)
		output += g.pop("rax")

		label := g.createLabel("if")
		output += "\ttest rax, rax\n"
		output += "\tjz " + label + "\n"
		output += g.GenScope(stmt.scope)
		output += label + ":\n"

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

		var variable Variable
		exists := false
		for _, v := range g.variables {
			if v.name == variableName {
				variable = v
				exists = true
			}
		}

		if !exists {
			panic(fmt.Errorf("unknown identifier: %v", variableName))
		}

		output += g.push(fmt.Sprintf("QWORD [rsp + %v]", (g.stackSize-variable.stackLoc-1)*8))

	case NodeTermRoundBracketExpr:
		output += g.GenExpr(term.expr)
	default:
		panic(fmt.Errorf("generator error: don't know how to generate term: %T", rawTerm.variant))
	}
	return output
}

func (g *Generator) GenScope(scope NodeScope) string {
	output := ""

	output += g.beginScope()

	for _, stmt := range scope.stmts {
		output += g.GenStmt(stmt)
	}

	output += g.endScope()

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

func (g *Generator) beginScope() string {
	g.scopes = append(g.scopes, len(g.variables))
	return ""
}
func (g *Generator) endScope() string {
	popCount := len(g.variables) - g.scopes[len(g.scopes)-1]

	g.stackSize -= uint(popCount)
	g.variables = g.variables[0 : len(g.variables)-popCount]
	g.scopes = g.scopes[0 : len(g.scopes)-1]

	return "\tadd rsp, " + fmt.Sprintf("%d", popCount*8) + "\n"
}

func (g *Generator) createLabel(labelCtx ...string) string {
	var suffix string
	if len(labelCtx) == 1 {
		suffix = "_" + labelCtx[0]
	} else {
		suffix = ""
	}

	g.labelCount++
	return fmt.Sprintf("label%d%s", g.labelCount, suffix)
}

type Variable struct {
	name     string
	stackLoc uint
}
