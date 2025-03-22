package main

import (
	"fmt"
	"slices"
	"strconv"
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

	inFunc          bool
	currentFunction Function

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
	g.inFunc = true
	output := ""

	functionName := stmt.ident.value.MustGetValue()

	returnCount, err := strconv.Atoi(stmt.returns)
	if err != nil {
		return "", err
	}

	for _, f := range g.functions {
		if f.name == functionName && f.returnCount == len(stmt.returns) {
			return "", stmt.ident.lineInfo.PositionedError(fmt.Sprintf("function identifier already used: %v", functionName))
		}
	}
	g.currentFunction = Function{
		name:        functionName,
		returnCount: returnCount,
	}

	typeSignature := ""
	paramTypes := []_Type{}
	for _, p := range stmt.params {
		t, err := _typeFromNodeType(p._type)
		if err != nil {
			return "", stmt.ident.lineInfo.PositionedError(err.Error())
		}

		paramTypes = append(paramTypes, t)
		typeSignature += t.getType() + "."
	}

	g.currentFunction.parameters = paramTypes

	signature := fmt.Sprintf("%s_%s", functionName, typeSignature)
	g.currentFunction.signature = signature

	output += fmt.Sprintf("%s:\n", signature)

	if g.genASMComments {
		output += "\t;=====FUNCTION SETUP=====\n"
	}
	output += g.push("rbp")
	output += "\tmov rbp, rsp\n\n"

	parameters := []Variable{}
	for i, p := range stmt.params {
		t, err := _typeFromNodeType(p._type)
		if err != nil {
			return "", stmt.ident.lineInfo.PositionedError(err.Error())
		}

		v := Variable{
			isParameter: true,
			name:        p.ident.value.MustGetValue(),
			stackLoc:    uint(i + 2),
			_type:       t,
		}
		parameters = append(parameters, v)
		g.stackSize++
	}

	body, err := g.GenFunctionBody(stmt.body, parameters)
	if err != nil {
		return "", err
	}
	output += body

	if g.genASMComments {
		output += "\t;=====FUNCTION CLEANUP=====\n"
	}

	output += g.pop("rbp")

	output += "\tret\n"

	g.inFunc = false
	g.functions = append(g.functions, g.currentFunction)
	g.currentFunction = Function{}
	g.stackSize = 0

	return output, nil
}

func (g *Generator) GenStmt(rawStmt NodeStmt) (string, error) {
	output := ""

	switch stmt := rawStmt.(type) {
	case NodeStmtVarDeclare:
		variableName := stmt.value.ident.value.MustGetValue()

		for _, v := range g.variables {
			if v.name == variableName {
				return "", stmt.value.ident.lineInfo.PositionedError(fmt.Sprintf("variable identifier already used: %v", variableName))
			}
		}

		t, err := _typeFromNodeType(stmt.value._type)
		if err != nil {
			return "", stmt.value.ident.lineInfo.PositionedError(err.Error())
		}

		g.variables = append(g.variables, Variable{
			stackLoc: g.stackSize,
			name:     variableName,
			_type:    t,
		})
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

		expr, err := g.GenIntExpr(stmt.expr)
		if err != nil {
			return "", err
		}
		output += expr
		output += g.pop("rax")
		if variable.isParameter {
			output += fmt.Sprintf("\tmov QWORD [rbp + %v], rax\n", variable.stackLoc*8)
		} else {
			output += fmt.Sprintf("\tmov QWORD [rsp + %v], rax\n", (g.stackSize-variable.stackLoc-1)*8)
		}

	case NodeStmtPointerAssign:
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

		expr, err := g.GenIntExpr(stmt.expr)
		if err != nil {
			return "", err
		}
		output += expr
		output += g.pop("rax")
		if variable.isParameter {
			output += fmt.Sprintf("\tmov rbx, [rbp + %v]\n", variable.stackLoc*8)
		} else {
			output += fmt.Sprintf("\tmov rbx, [rsp + %v]\n", (g.stackSize-variable.stackLoc-1)*8)
		}
		output += "\tmov QWORD [rbx], rax\n"

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

		expr, err := g.GenIntExpr(stmt.expr)
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

	case NodeFunctionCall:
		funcCall, retCount, err := g.GenFuncCall(stmt)
		if err != nil {
			return "", err
		}
		output += funcCall

		// get rid of the return values as we're not storing them
		output += "\tadd rsp, " + fmt.Sprintf("%d", retCount*8) + "\n"
		g.stackSize -= uint(retCount)

	case NodeStmtReturn:
		if !g.inFunc {
			return "", stmt._return.lineInfo.PositionedError("can only return when in a function")
		}

		if len(stmt.returns) != g.currentFunction.returnCount {
			return "", stmt._return.lineInfo.PositionedError(fmt.Sprintf("incorrect number of values returned. Expected %v, Found %v", g.currentFunction.returnCount, len(stmt.returns)))
		}

		for i, expr := range stmt.returns {
			expr, err := g.GenIntExpr(expr)
			if err != nil {
				return "", err
			}
			output += expr

			stackOffset := (len(g.currentFunction.parameters) + i + 2) * 8

			output += g.pop(fmt.Sprintf("QWORD [rbp + %v]", stackOffset))
		}

		output += g.exitFunction(false)
		output += "\tpop rbp\n" //g.pop("rbp")
		output += "\tret\n"

	case NodeStmtSyscall:
		argRegisters := []string{"rax", "rdi", "rsi", "rdx", "r10", "r8", "r9"}
		usedArgs := len(stmt.arguments)

		for _, e := range stmt.arguments {
			expr, err := g.GenIntExpr(e)
			if err != nil {
				return "", err
			}
			output += expr
		}

		for i := usedArgs - 1; i >= 0; i-- {
			output += g.pop(argRegisters[i])
		}
		output += "\tsyscall\n"

	default:
		panic(fmt.Errorf("generator error: don't know how to generate statement: %T", rawStmt))
	}
	return output, nil
}

func (g *Generator) GenFuncCall(stmt NodeFunctionCall) (string, int, error) {
	output := ""

	functionName := stmt.ident.value.MustGetValue()
	var function Function
	exists := false
	foundWrong := false

	allFunctions := append(g.functions, g.currentFunction)
	for _, f := range allFunctions {
		if f.name == functionName {
			exists = true
			if len(stmt.params) == len(f.parameters) {
				function = f
				foundWrong = false
				break
			}
			foundWrong = true
		}
	}
	if !exists {
		return "", 0, stmt.ident.lineInfo.PositionedError(fmt.Sprintf("undefined function: '%s'", functionName))
	}
	if foundWrong {
		return "", 0, stmt.ident.lineInfo.PositionedError("incorrect number of arguments passed.")
	}

	for i := 0; i < function.returnCount; i++ {
		output += g.push("0")
	}

	reversed := stmt.params
	slices.Reverse(reversed)
	for _, p := range reversed {
		expr, err := g.GenIntExpr(p)
		if err != nil {
			return "", 0, err
		}
		output += expr
	}

	output += fmt.Sprintf("\tcall %s\n", function.signature)

	output += "\tadd rsp, " + fmt.Sprintf("%d", len(stmt.params)*8) + "\n"
	g.stackSize -= uint(len(stmt.params))

	// Must deal with cleaning up return values from go-stack at function call-site

	return output, function.returnCount, nil
}

func (g *Generator) GenIntExpr(rawExpr NodeExpr) (string, error) {
	output := ""

	switch expr := rawExpr.(type) {
	case NodeIntTerm:
		term, err := g.GenIntTerm(expr)
		if err != nil {
			return "", err
		}
		output += term
	case NodeIntBinExpr:
		binExpr, err := g.GenIntBinExpr(expr)
		if err != nil {
			return "", err
		}
		output += binExpr
	default:
		panic(fmt.Errorf("generator error: don't know how to generate expression: %T", rawExpr))
	}
	return output, nil
}

func (g *Generator) GenIntBinExpr(rawBinExpr NodeIntBinExpr) (string, error) {
	output := ""
	switch binExpr := rawBinExpr.(type) {
	case NodeIntBinExprAdd:
		expr, err := g.GenIntExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenIntExpr(binExpr.right)
		if err != nil {
			return "", err
		}

		output += expr
		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tadd rax, rbx\n"
		output += g.push("rax")
	case NodeIntBinExprSubtract:
		expr, err := g.GenIntExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenIntExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tsub rax, rbx\n"
		output += g.push("rax")
	case NodeIntBinExprMultiply:
		expr, err := g.GenIntExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenIntExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tmul rbx\n"
		output += g.push("rax")
	case NodeIntBinExprDivide:
		expr, err := g.GenIntExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenIntExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tmov rdx, 0\n"
		output += "\tdiv rbx\n"
		output += g.push("rax")
	case NodeIntBinExprModulo:
		expr, err := g.GenIntExpr(binExpr.left)
		if err != nil {
			return "", err
		}
		output += expr
		expr, err = g.GenIntExpr(binExpr.right)
		if err != nil {
			return "", err
		}
		output += expr

		output += g.pop("rbx")
		output += g.pop("rax")
		output += "\tmov rdx, 0\n"
		output += "\tdiv rbx\n"
		output += g.push("rdx")
	default:
		panic(fmt.Errorf("generator error: don't know how to generate binary expression: %T", rawBinExpr))
	}
	return output, nil
}

func (g *Generator) GenIntTerm(rawTerm NodeIntTerm) (string, error) {
	output := ""

	switch term := rawTerm.(type) {
	case NodeIntTermNegativeTerm:
		subTerm, err := g.GenIntTerm(term.term)
		if err != nil {
			return "", err
		}
		output += subTerm

		output += g.pop("rax")
		output += "\tneg rax\n"
		output += g.push("rax")
	case NodeIntTermLiteral:
		output += "\tmov rax, " + term.intLiteral.value.MustGetValue() + "\n"
		output += g.push("rax")

	case NodeIntTermIdentifier:
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

		if variable.isParameter {
			output += g.push(fmt.Sprintf("QWORD [rbp + %v]", variable.stackLoc*8))
		} else {
			output += g.push(fmt.Sprintf("QWORD [rsp + %v]", (g.stackSize-variable.stackLoc-1)*8))
		}

	case NodeFunctionCall:
		funcCall, retCount, err := g.GenFuncCall(term)
		if err != nil {
			return "", err
		}

		if retCount != 1 {
			return "", term.ident.lineInfo.PositionedError("function doesn't return any values (or more than 1 atm) so can't be used as a term")
		}

		output += funcCall

	case NodeIntTermRoundBracketExpr:
		expr, err := g.GenIntExpr(term.expr)
		if err != nil {
			return "", err
		}
		output += expr

	case NodeIntTermPointer:
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

		if variable.isParameter {
			output += "\tlea rax, QWORD [rbp + " + fmt.Sprint(variable.stackLoc*8) + "]\n"
		} else {
			output += "\tlea rax, QWORD [rsp + " + fmt.Sprint((g.stackSize-variable.stackLoc-1)*8) + "]\n"
		}
		output += g.push("rax")

	case NodeIntTermPointerDereference:
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

		if variable.isParameter {
			output += "\tmov rax, [rbp + " + fmt.Sprint(variable.stackLoc*8) + "]\n"
		} else {
			output += "\tmov rax, [rsp + " + fmt.Sprint((g.stackSize-variable.stackLoc-1)*8) + "]\n"
		}

		output += "\tmov rax, [rax]\n"

		output += g.push("rax")

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

func (g *Generator) GenFunctionBody(body NodeScope, params []Variable) (string, error) {
	output := ""

	// start "scope" for the function body
	g.scopes = append(g.scopes, len(g.variables))

	g.currentFunction.scopeIndex = len(g.scopes) - 1

	g.variables = append(g.variables, params...)

	if g.genASMComments {
		output += "\t;=====FUNCTION BODY=====\n"
	}
	for _, stmt := range body.stmts {
		generated, err := g.GenStmt(stmt)
		if err != nil {
			return "", err
		}
		output += generated + "\n"
	}

	output += g.exitFunction(true)

	return output, nil
}

func (g *Generator) GenIf(_if NodeStmtIf) (string, error) {
	output := ""

	expr, err := g.GenIntExpr(_if.expr)
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
	output := ""

	g.scopes = append(g.scopes, len(g.variables))

	if g.genASMComments {
		output += "\t;---start_scope---\n"
	}

	return output
}
func (g *Generator) endScope() string {
	output := ""

	popCount := len(g.variables) - g.scopes[len(g.scopes)-1]

	g.stackSize -= uint(popCount)
	g.variables = g.variables[0 : len(g.variables)-popCount]
	g.scopes = g.scopes[0 : len(g.scopes)-1]

	output += "\tadd rsp, " + fmt.Sprintf("%d", popCount*8) + "\n"

	if g.genASMComments {
		output += "\t;---end_scope---\n"
	}

	return output
}
func (g *Generator) exitFunction(updateLocal bool) string {
	targetVariableCount := g.scopes[g.currentFunction.scopeIndex] + len(g.currentFunction.parameters)
	popCount := len(g.variables) - targetVariableCount

	// localPopCount := len(g.variables) - g.scopes[len(g.scopes)-1]

	// for i := len(g.scopes) - 1; i >= g.currentFunction.scopeIndex; i-- {
	// 	len(g.variables)- g.scopes[]
	// }
	// popCount := localPopCount - paramCount

	if updateLocal {
		localPopCount := popCount + len(g.currentFunction.parameters)

		g.stackSize -= uint(localPopCount)
		g.variables = g.variables[0 : len(g.variables)-localPopCount]
		g.scopes = g.scopes[0:g.currentFunction.scopeIndex]
	}

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

type _Type interface {
	getType() string
}

func _typeFromNodeType(n NodeType) (_Type, error) {
	switch n := n.(type) {
	case NodePointerType:
		sub, err := _typeFromNodeType(n.subType)
		if err != nil {
			return nil, err
		}

		return _PointerType{
			subType: sub,
		}, nil
	case NodePureType:
		return _typeFromNodeType(n.baseType)
	case NodeBoolType:
		return boolType{}, nil
	case NodeIntType:
		return intType{}, nil
	case NodeCharType:
		return charType{}, nil
	default:
		return nil, fmt.Errorf("unknown type: %T", n)
	}

}

type boolType struct{}

func (boolType) getType() string {
	return "B"
}

type intType struct{}

func (intType) getType() string {
	return "I"
}

type charType struct{}

func (charType) getType() string {
	return "C"
}

type _PointerType struct {
	subType _Type
}

func (p _PointerType) getType() string {
	return "P" + p.subType.getType()
}

type Variable struct {
	name     string
	stackLoc uint

	isParameter bool
	_type       _Type
}

type Function struct {
	name        string
	signature   string
	returnCount int
	parameters  []_Type

	scopeIndex int
}
