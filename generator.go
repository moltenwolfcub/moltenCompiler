package main

import (
	"fmt"
	"slices"
)

type Generator struct {
	program NodeProg

	stackSize uint
	variables []Variable
	functions []Function
	scopes    []int

	labelCount    int
	breakLabel    string
	continueLabel string

	genASMComments bool
}

func NewGenerator(prog NodeProg) Generator {
	return Generator{
		program: prog,

		stackSize: 0,
		variables: []Variable{},
		functions: []Function{},
		scopes:    []int{},

		labelCount:    0,
		breakLabel:    "nil",
		continueLabel: "nil",

		genASMComments: true,
	}
}

func (g *Generator) GenProg() (string, error) {
	output := "global _start\n\n\n"

	pre, err := g.PreGenerate()
	if err != nil {
		return "", err
	}
	output += pre

	output += "_start:\n"

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

func (g *Generator) PreGenerate() (string, error) {
	output := ""

	for _, stmt := range g.program.stmts {
		funcStmt, ok := stmt.(NodeStmtFunctionDefinition)
		if !ok {
			continue
		}
		generated, err := g.GenFuncDefinition(funcStmt)
		if err != nil {
			return "", err
		}
		output += generated + "\n\n"
	}
	return output, nil
}

func (g *Generator) GenFuncDefinition(stmt NodeStmtFunctionDefinition) (string, error) {
	output := ""

	functionName := stmt.ident.value.MustGetValue()

	for _, f := range g.functions {
		if f.name == functionName {
			return "", stmt.ident.lineInfo.PositionedError(fmt.Sprintf("function identifier already used: %v", functionName))
		}
	}
	g.functions = append(g.functions, Function{name: functionName})

	output += functionName + ":\n"

	parameters := []Variable{}
	for _, p := range stmt.params {
		v := Variable{
			name: p.value.MustGetValue(),
		}
		parameters = append(parameters, v)
	}

	body, err := g.GenScopeWithParams(stmt.body, parameters)
	if err != nil {
		return "", err
	}
	output += body

	output += "\tret\n"

	return output, nil
}

func (g *Generator) GenStmt(rawStmt NodeStmt) (string, error) {
	output := ""

	switch stmt := rawStmt.(type) {
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

		for _, v := range g.variables {
			if v.name == variableName {
				return "", stmt.ident.lineInfo.PositionedError(fmt.Sprintf("variable identifier already used: %v", variableName))
			}
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
				break
			}
		}
		if !exists {
			return "", stmt.ident.lineInfo.PositionedError(fmt.Sprintf("undefined variable: '%s'", variableName))
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
			return "", stmt._break.lineInfo.PositionedError("can't break when not in a loop")
		}
		output += "\tjmp " + g.breakLabel + "\n"

	case NodeStmtContinue:
		if g.continueLabel == "nil" {
			return "", stmt._continue.lineInfo.PositionedError("can't continue when not in a loop")
		}
		output += "\tjmp " + g.continueLabel + "\n"

	case NodeStmtFunctionDefinition:
		/*
			functions are generated before other statements
			so they are in the correct order.

			this is just here to keep the type switch happy.
		*/

	case NodeStmtFunctionCall:
		functionName := stmt.ident.value.MustGetValue()
		var function Function
		exists := false
		for _, f := range g.functions {
			if f.name == functionName {
				function = f
				exists = true
				break
			}
		}
		if !exists {
			return "", stmt.ident.lineInfo.PositionedError(fmt.Sprintf("undefined function: '%s'", functionName))
		}

		reversed := stmt.params
		slices.Reverse(reversed)
		for _, p := range reversed {
			expr, err := g.GenExpr(p)
			if err != nil {
				return "", err
			}
			output += expr
		}

		output += "\tcall " + function.name + "\n"

	default:
		panic(fmt.Errorf("generator error: don't know how to generate statement: %T", rawStmt))
	}
	return output, nil
}

func (g *Generator) GenExpr(rawExpr NodeExpr) (string, error) {
	output := ""

	switch expr := rawExpr.(type) {
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
		panic(fmt.Errorf("generator error: don't know how to generate expression: %T", rawExpr))
	}
	return output, nil
}

func (g *Generator) GenBinExpr(rawBinExpr NodeBinExpr) (string, error) {
	output := ""
	switch binExpr := rawBinExpr.(type) {
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
		panic(fmt.Errorf("generator error: don't know how to generate binary expression: %T", rawBinExpr))
	}
	return output, nil
}

func (g *Generator) GenTerm(rawTerm NodeTerm) (string, error) {
	output := ""

	switch term := rawTerm.(type) {
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
			return "", term.identifier.lineInfo.PositionedError(fmt.Sprintf("undefined variable: %v", variableName))
		}

		output += g.push(fmt.Sprintf("QWORD [rsp + %v]", (g.stackSize-variable.stackLoc-1)*8))

	case NodeTermRoundBracketExpr:
		expr, err := g.GenExpr(term.expr)
		if err != nil {
			return "", err
		}
		output += expr
	default:
		panic(fmt.Errorf("generator error: don't know how to generate term: %T", rawTerm))
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
		output += generated + "\n"
	}

	output += g.endScope()

	return output, nil
}

func (g *Generator) GenScopeWithParams(scope NodeScope, params []Variable) (string, error) {
	output := ""

	output += g.beginScope()

	if len(params) > 0 && g.genASMComments {
		output += "\t;=====PARAMETERS=====\n"
	}
	for i, p := range params {
		g.variables = append(g.variables, Variable{stackLoc: g.stackSize, name: p.name})

		stackOffset := 8 + (i)*16

		if g.genASMComments {
			output += "\t;" + p.name + "\n"
		}
		/*	start at offset 8 and jump by 16 to account for the other parameters
			that have just been pushed to the stack.
		*/
		output += fmt.Sprintf("\tmov rax, QWORD [rsp + %v]\n", stackOffset)
		output += g.push("rax")
		output += "\n"
	}
	if len(params) > 0 {
		output += "\n"
	}

	if g.genASMComments {
		output += "\t;=====FUNCTION BODY=====\n"
	}
	for _, stmt := range scope.stmts {
		generated, err := g.GenStmt(stmt)
		if err != nil {
			return "", err
		}
		output += generated + "\n"
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

	switch _else := rawElse.(type) {
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
		panic(fmt.Errorf("generator error: don't know how to generate else branch: %T", rawElse))
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

type Function struct {
	name string
}
