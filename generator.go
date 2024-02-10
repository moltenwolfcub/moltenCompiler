package main

import (
	"fmt"
)

type Generator struct {
	program       NodeProg
	stackSize     uint
	variables     []Variable
	scopes        []int
	labelCount    int
	breakLabel    string
	continueLabel string
}

func NewGenerator(prog NodeProg) Generator {
	return Generator{
		program:       prog,
		stackSize:     0,
		variables:     []Variable{},
		scopes:        []int{},
		labelCount:    0,
		breakLabel:    "nil",
		continueLabel: "nil",
	}
}

func (g *Generator) GenProg() (string, error) {
	output := "global _start\n_start:\n"

	for _, stmt := range g.program.stmts {
		generated, err := g.GenStmt(stmt)
		if err != nil {
			return "", err
		}
		output += generated + "\n"
	}

	//exit 0 at end of program if no explicit exit called
	output += "\tmov rax, 60\n"
	output += "\tmov rdi, 0\n"
	output += "\tsyscall\n"

	return output, nil
}

func (g *Generator) GenStmt(rawStmt NodeStmt) (string, error) {
	output := ""

	switch stmt := rawStmt.variant.(type) {
	case NodeStmtExit:
		expr, err := g.GenExpr(stmt.expr)
		if err != nil {
			return "", err
		}
		output += expr
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
			return "", g.error(stmt.ident, fmt.Sprintf("identifier already used: %v", variableName))
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
			return "", g.error(stmt.ident, fmt.Sprintf("variables must be declared before assignment. '%s' is undefined", variableName))
		}

		expr, err := g.GenExpr(stmt.expr)
		if err != nil {
			return "", err
		}
		output += expr
		output += g.pop("rax")
		output += fmt.Sprintf("\tmov QWORD [rsp + %v], rax\n", (g.stackSize-variable.stackLoc-1)*8)

	case NodeScope:
		scope, err := g.GenScope(stmt)
		if err != nil {
			return "", err
		}
		output += scope
	case NodeStmtIf:
		ifStmt, err := g.GenIf(stmt)
		if err != nil {
			return "", err
		}
		output += ifStmt

	case NodeStmtWhile:
		startLabel := g.createLabel("startWhile")
		endLabel := g.createLabel("endWhile")

		g.breakLabel = endLabel
		g.continueLabel = startLabel

		output += startLabel + ":\n"

		expr, err := g.GenExpr(stmt.expr)
		if err != nil {
			return "", err
		}
		output += expr
		output += g.pop("rax")

		output += "\ttest rax, rax\n"
		output += "\tjz " + endLabel + "\n"

		scope, err := g.GenScope(stmt.scope)
		if err != nil {
			return "", err
		}
		output += scope

		output += "\tjmp " + startLabel + "\n"

		output += endLabel + ":\n"

		g.breakLabel = "nil"
		g.continueLabel = "nil"

	case NodeStmtBreak:
		if g.breakLabel == "nil" {
			return "", g.error(stmt._break, "can't break when not in a loop")
		}
		output += "\tjmp " + g.breakLabel + "\n"

	case NodeStmtContinue:
		if g.continueLabel == "nil" {
			return "", g.error(stmt._continue, "can't continue when not in a loop")
		}
		output += "\tjmp " + g.continueLabel + "\n"

	default:
		panic(fmt.Errorf("generator error: don't know how to generate statement: %T", rawStmt.variant))
	}
	return output, nil
}

func (g *Generator) GenExpr(rawExpr NodeExpr) (string, error) {
	output := ""

	switch expr := rawExpr.variant.(type) {
	case NodeTerm:
		term, err := g.GenTerm(expr)
		if err != nil {
			return "", err
		}
		output += term
	case NodeBinExpr:
		binExpr, err := g.GenBinExpr(expr)
		if err != nil {
			return "", err
		}
		output += binExpr
	default:
		panic(fmt.Errorf("generator error: don't know how to generate expression: %T", rawExpr.variant))
	}
	return output, nil
}

func (g *Generator) GenBinExpr(rawBinExpr NodeBinExpr) (string, error) {
	output := ""
	switch binExpr := rawBinExpr.variant.(type) {
	case NodeBinExprAdd:
		expr, err := g.GenExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenExpr(binExpr.right)
		if err != nil {
			return "", err
		}

		output += expr
		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tadd rax, rbx\n"
		output += g.push("rax")
	case NodeBinExprSubtract:
		expr, err := g.GenExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tsub rax, rbx\n"
		output += g.push("rax")
	case NodeBinExprMultiply:
		expr, err := g.GenExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tmul rbx\n"
		output += g.push("rax")
	case NodeBinExprDivide:
		expr, err := g.GenExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tdiv rbx\n"
		output += g.push("rax")
	default:
		panic(fmt.Errorf("generator error: don't know how to generate binary expression: %T", rawBinExpr.variant))
	}
	return output, nil
}

func (g *Generator) GenTerm(rawTerm NodeTerm) (string, error) {
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
			return "", g.error(term.identifier, fmt.Sprintf("unknown identifier: %v", variableName))
		}

		output += g.push(fmt.Sprintf("QWORD [rsp + %v]", (g.stackSize-variable.stackLoc-1)*8))

	case NodeTermRoundBracketExpr:
		expr, err := g.GenExpr(term.expr)
		if err != nil {
			return "", err
		}
		output += expr
	default:
		panic(fmt.Errorf("generator error: don't know how to generate term: %T", rawTerm.variant))
	}
	return output, nil
}

func (g *Generator) GenScope(scope NodeScope) (string, error) {
	output := ""

	output += g.beginScope()

	for _, stmt := range scope.stmts {
		generated, err := g.GenStmt(stmt)
		if err != nil {
			return "", err
		}
		output += generated
	}

	output += g.endScope()

	return output, nil
}

func (g *Generator) GenIf(_if NodeStmtIf) (string, error) {
	output := ""

	expr, err := g.GenExpr(_if.expr)
	if err != nil {
		return "", err
	}
	output += expr

	output += g.pop("rax")

	label := g.createLabel("else")
	output += "\ttest rax, rax\n"
	output += "\tjz " + label + "\n"
	scope, err := g.GenScope(_if.scope)
	if err != nil {
		return "", err
	}
	output += scope
	output += label + ":\n"

	if _if.elseBranch.HasValue() {
		_else, err := g.GenElse(_if.elseBranch.MustGetValue())
		if err != nil {
			return "", err
		}
		output += _else
	}

	return output, nil
}

func (g *Generator) GenElse(rawElse NodeElse) (string, error) {
	output := ""

	switch _else := rawElse.variant.(type) {
	case NodeElseScope:
		scope, err := g.GenScope(_else.scope)
		if err != nil {
			return "", err
		}
		output += scope
	case NodeElseElif:
		ifStmt, err := g.GenIf(_else.ifStmt)
		if err != nil {
			return "", err
		}
		output += ifStmt
	default:
		panic(fmt.Errorf("generator error: don't know how to generate else branch: %T", rawElse.variant))
	}
	return output, nil

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

func (g Generator) error(t Token, message string) error {
	return fmt.Errorf("%s:%d:%d: %s", t.file, t.line, t.col, message)
}
