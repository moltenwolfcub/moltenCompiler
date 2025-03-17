package main

import (
	"errors"

	opt "github.com/moltenwolfcub/moltenCompiler/optional"
)

type Parser struct {
	tokens       []Token
	currentIndex int
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) ParseProg() (NodeProg, error) {
	node := NodeProg{
		[]NodeStmt{},
	}
	for p.peek().HasValue() {
		stmt, err := p.ParseStmt()
		if err != nil {
			return NodeProg{}, p.error(err.Error())
		}

		node.stmts = append(node.stmts, stmt)
	}
	return node, nil
}

func (p *Parser) ParseStmt() (NodeStmt, error) {
	if p.mustTryConsume(_var).HasValue() {
		node := NodeStmtVarDeclare{}

		typedIdent := TypedIdentifier{}

		tok, err := p.tryConsume(identifier, "expected variable identifier after `var`")
		if err != nil {
			return nil, err
		}

		typedIdent.ident = tok

		_type, err := p.ParseType()
		if err != nil {
			return nil, err
		}

		typedIdent._type = _type

		node.value = typedIdent

		_, err = p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return nil, err
		}

		return node, nil
	} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == identifier {
		if !p.peek(1).HasValue() {
			return nil, errors.New("expected '=' or '()' after identifier for variable assignment or function call. didn't find any token")
		}

		switch p.peek(1).MustGetValue().tokenType {
		case equals:

			node := NodeStmtVarAssign{
				ident: p.consume(),
			}
			p.consume()

			expr, err := p.ParseIntExpr()
			if err != nil {
				return nil, err
			}
			node.expr = expr

			_, err = p.tryConsume(semiColon, "missing ';'")
			if err != nil {
				return nil, err
			}

			return node, nil
		case openRoundBracket:
			funcCall, err := p.ParseFuncCall()
			if err != nil {
				return nil, err
			}

			_, err = p.tryConsume(semiColon, "missing ';'")
			if err != nil {
				return nil, err
			}

			return funcCall, nil
		default:
			return nil, errors.New("expected '=' or '()' after identifier for variable assignment or function call")
		}
	} else if p.mustTryConsume(asterisk).HasValue() {
		tok, err := p.tryConsume(identifier, "expected variable identifier after '*'")
		if err != nil {
			return nil, err
		}

		node := NodeStmtPointerAssign{
			ident: tok,
		}

		_, err = p.tryConsume(equals, "expected '=' after identifier for pointer assignment")
		if err != nil {
			return nil, err
		}

		expr, err := p.ParseIntExpr()
		if err != nil {
			return nil, err
		}
		node.expr = expr

		_, err = p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return nil, err
		}

		return node, nil

	} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == openCurlyBracket {
		scope, err := p.ParseScope()
		if err != nil {
			return nil, err
		}
		return scope, nil

	} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == _if {
		ifStmt, err := p.ParseIf()
		if err != nil {
			return nil, err
		}
		return ifStmt, nil

	} else if p.mustTryConsume(while).HasValue() {
		node := NodeStmtWhile{}

		_, err := p.tryConsume(openRoundBracket, "Expected '('")
		if err != nil {
			return nil, err
		}

		node.expr, err = p.ParseIntExpr()
		if err != nil {
			return nil, err
		}

		_, err = p.tryConsume(closeRoundBracket, "Expected ')'")
		if err != nil {
			return nil, err
		}

		node.scope, err = p.ParseScope()
		if err != nil {
			return nil, err
		}
		return node, nil

	} else if tok := p.mustTryConsume(_break); tok.HasValue() {
		_, err := p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return nil, err
		}
		return NodeStmtBreak{tok.MustGetValue()}, nil

	} else if tok := p.mustTryConsume(_continue); tok.HasValue() {
		_, err := p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return nil, err
		}
		return NodeStmtContinue{tok.MustGetValue()}, nil

	} else if p.mustTryConsume(_func).HasValue() {
		node := NodeStmtFunctionDefinition{}

		count, err := p.tryConsume(intLiteral, "expected an int for the number of returns")
		if err != nil {
			return nil, err
		}
		node.returns = count.value.MustGetValue()

		ident, err := p.tryConsume(identifier, "expected function identifier after `func`")
		if err != nil {
			return nil, err
		}
		node.ident = ident

		_, err = p.tryConsume(openRoundBracket, "missing '('")
		if err != nil {
			return nil, err
		}

		for {
			ident, err := p.tryConsume(identifier, "optional so this should never error")
			if err != nil {
				break
			}
			_type, err := p.ParseType()
			if err != nil {
				return nil, err
			}

			param := TypedIdentifier{
				ident: ident,
				_type: _type,
			}

			node.params = append(node.params, param)

			_, err = p.tryConsume(comma, "optional so this should never error")
			if err != nil {
				break
			}
		}

		_, err = p.tryConsume(closeRoundBracket, "missing ')'")
		if err != nil {
			return nil, err
		}

		scope, err := p.ParseScope()
		if err != nil {
			return nil, err
		}
		node.body = scope

		return node, nil
	} else if tok := p.mustTryConsume(_return); tok.HasValue() {
		node := NodeStmtReturn{_return: tok.MustGetValue()}

		for {
			expr, err := p.ParseIntExpr()
			if err == errMissingExpr {
				break
			} else if err != nil {
				return NodeIntFunctionCall{}, err
			}
			node.returns = append(node.returns, expr)

			_, err = p.tryConsume(comma, "optional so this should never error")
			if err != nil {
				break
			}
		}

		_, err := p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return nil, err
		}

		return node, nil

	} else if tok := p.mustTryConsume(syscall); tok.HasValue() {
		node := NodeStmtSyscall{syscall: tok.MustGetValue()}

		_, err := p.tryConsume(openRoundBracket, "Expected '('")
		if err != nil {
			return nil, err
		}

		for {
			expr, err := p.ParseIntExpr()
			if err == errMissingExpr {
				break
			} else if err != nil {
				return NodeStmtSyscall{}, err
			}
			node.arguments = append(node.arguments, expr)

			_, err = p.tryConsume(comma, "optional so this should never error")
			if err != nil {
				break
			}
		}
		if len(node.arguments) > 7 {
			return nil, errSyscallArgs
		}

		_, err = p.tryConsume(closeRoundBracket, "Expected ')'")
		if err != nil {
			return nil, err
		}

		_, err = p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return nil, err
		}

		return node, nil

	} else {
		return nil, errMissingStmt
	}
}

var errSyscallArgs error = errors.New("syscalls can't have more than 7 arguments")
var errMissingStmt error = errors.New("expected statement but couldn't find one")

// based off of this principle and algorithm:
// https://eli.thegreenplace.net/2012/08/02/parsing-expressions-by-precedence-climbing
func (p *Parser) ParseIntExpr(minPrecedence ...int) (NodeExpr, error) {
	var minPrec int
	if len(minPrecedence) == 1 {
		minPrec = minPrecedence[0]
	} else {
		minPrec = 0
	}

	lhsTerm, err := p.ParseIntTerm()
	if err == errMissingTerm {
		return nil, errMissingExpr
	}
	if err != nil {
		return nil, err
	}

	lhsExpr := NodeExpr(lhsTerm)

	for {
		currentToken := p.peek()
		if !currentToken.HasValue() {
			break
		}
		currentPrec := currentToken.MustGetValue().tokenType.GetBinPrec()
		if !currentPrec.HasValue() || currentPrec.MustGetValue() < minPrec { //prolly meant to be <=
			break
		}
		op := p.consume()

		nextMinPrec := currentPrec.MustGetValue() + 1

		rhsExpr, err := p.ParseIntExpr(nextMinPrec)
		if err != nil {
			return nil, err
		}

		var expr NodeIntBinExpr
		switch op.tokenType {
		case plus:
			expr = NodeIntBinExprAdd{
				left:  lhsExpr,
				right: rhsExpr,
			}
		case asterisk:
			expr = NodeIntBinExprMultiply{
				left:  lhsExpr,
				right: rhsExpr,
			}
		case minus:
			expr = NodeIntBinExprSubtract{
				left:  lhsExpr,
				right: rhsExpr,
			}
		case fslash:
			expr = NodeIntBinExprDivide{
				left:  lhsExpr,
				right: rhsExpr,
			}
		case percent:
			expr = NodeIntBinExprModulo{
				left:  lhsExpr,
				right: rhsExpr,
			}
		}
		lhsExpr = expr

	}
	return lhsExpr, nil
}

var errMissingExpr error = errors.New("expected expression")

func (p *Parser) ParseIntTerm() (NodeIntTerm, error) {
	if p.mustTryConsume(minus).HasValue() {
		term, err := p.ParseIntTerm()
		if err != nil {
			return nil, err
		}
		return NodeIntTermNegativeTerm{term}, nil
	} else if tok := p.mustTryConsume(intLiteral); tok.HasValue() {
		return NodeIntTermLiteral{tok.MustGetValue()}, nil
	} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == identifier {
		if p.peek(1).MustGetValue().tokenType == openRoundBracket {
			return p.ParseFuncCall()
		} else {
			return NodeIntTermIdentifier{p.consume()}, nil
		}
	} else if p.mustTryConsume(openRoundBracket).HasValue() {
		expr, err := p.ParseIntExpr()
		if err != nil {
			return nil, err
		}
		_, err = p.tryConsume(closeRoundBracket, "expected ')'")
		if err != nil {
			return nil, err
		}
		return NodeIntTermRoundBracketExpr{expr}, nil
	} else if p.mustTryConsume(ampersand).HasValue() {
		variable, err := p.tryConsume(identifier, "expected variable identifier after '&'")
		if err != nil {
			return nil, err
		}
		return NodeIntTermPointer{variable}, nil
	} else if p.mustTryConsume(asterisk).HasValue() {
		variable, err := p.tryConsume(identifier, "expected variable identifier after '*'")
		if err != nil {
			return nil, err
		}
		return NodeIntTermPointerDereference{variable}, nil
	}
	return nil, errMissingTerm
}

var errMissingTerm error = errors.New("expected term but couldn't find one")

func (p *Parser) ParseScope() (NodeScope, error) {
	if !p.mustTryConsume(openCurlyBracket).HasValue() {
		return NodeScope{}, nil
	}

	var scope NodeScope
	for {
		stmt, err := p.ParseStmt()
		if err == errMissingStmt {
			break
		} else if err != nil {
			return NodeScope{}, err
		}

		scope.stmts = append(scope.stmts, stmt)
	}
	_, err := p.tryConsume(closeCurlyBracket, "expected '}'")
	if err != nil {
		return NodeScope{}, err
	}

	return scope, nil
}

var errMissingScopeStmt error = errors.New("expected a scope statement but didn't find an open brace")

func (p *Parser) ParseIf() (NodeStmtIf, error) {
	if !p.mustTryConsume(_if).HasValue() {
		return NodeStmtIf{}, errMissingIfStmt
	}

	node := NodeStmtIf{}

	_, err := p.tryConsume(openRoundBracket, "Expected '('")
	if err != nil {
		return NodeStmtIf{}, err
	}

	node.expr, err = p.ParseIntExpr()
	if err != nil {
		return NodeStmtIf{}, err
	}

	_, err = p.tryConsume(closeRoundBracket, "Expected ')'")
	if err != nil {
		return NodeStmtIf{}, err
	}

	node.scope, err = p.ParseScope()
	if err != nil {
		return NodeStmtIf{}, err
	}

	node.elseBranch, err = p.ParseElse()
	if err != nil {
		return NodeStmtIf{}, err
	}

	return node, nil
}

var errMissingIfStmt error = errors.New("expected if statement but didn't find `if` token")

func (p *Parser) ParseElse() (opt.Optional[NodeElse], error) {
	if !p.mustTryConsume(_else).HasValue() {
		return opt.Optional[NodeElse]{}, nil
	}

	var node NodeElse

	ifStmt, err := p.ParseIf()

	if err == errMissingIfStmt {
		scope, err := p.ParseScope()
		if err == errMissingScopeStmt {
			return opt.Optional[NodeElse]{}, errors.New("expected scope or else-if statement following `else` keyword")

		} else if err != nil {
			return opt.Optional[NodeElse]{}, err
		}
		node = NodeElseScope{scope}
		return opt.ToOptional(node), nil

	} else if err != nil {
		return opt.Optional[NodeElse]{}, err

	} else {
		node = NodeElseElif{ifStmt}
		return opt.ToOptional(node), nil
	}
}

func (p *Parser) ParseFuncCall() (NodeIntFunctionCall, error) {

	node := NodeIntFunctionCall{
		ident: p.consume(),
	}

	_, err := p.tryConsume(openRoundBracket, "missing '('")
	if err != nil {
		return NodeIntFunctionCall{}, err
	}

	for {
		expr, err := p.ParseIntExpr()
		if err == errMissingExpr {
			break
		} else if err != nil {
			return NodeIntFunctionCall{}, err
		}
		node.params = append(node.params, expr)

		_, err = p.tryConsume(comma, "optional so this should never error")
		if err != nil {
			break
		}
	}

	_, err = p.tryConsume(closeRoundBracket, "missing ')'")
	if err != nil {
		return NodeIntFunctionCall{}, err
	}

	return node, nil
}

func (p *Parser) ParseType() (NodeType, error) {
	if tok := p.mustTryConsume(asterisk); tok.HasValue() {
		_type, err := p.ParseType()
		if err != nil {
			return nil, err
		}
		return NodePointerType{_type}, nil
	} else {
		baseType, err := p.ParseBaseType()
		if err != nil {
			return nil, err
		}
		return NodePureType{baseType}, nil
	}
}

func (p *Parser) ParseBaseType() (NodeBaseType, error) {
	if tok := p.mustTryConsume(typeBool); tok.HasValue() {
		return NodeBoolType{_bool: tok.MustGetValue()}, nil
	} else if tok := p.mustTryConsume(typeInt); tok.HasValue() {
		return NodeIntType{_int: tok.MustGetValue()}, nil
	} else if tok := p.mustTryConsume(typeChar); tok.HasValue() {
		return NodeCharType{_char: tok.MustGetValue()}, nil
	}
	return nil, errMissingBaseType
}

var errMissingBaseType error = errors.New("expected base type")

func (p Parser) peek(offset ...int) opt.Optional[Token] {
	var offsetAmount int
	if len(offset) == 1 {
		offsetAmount = offset[0]
	} else {
		offsetAmount = 0
	}

	if p.currentIndex+offsetAmount >= len(p.tokens) {
		return opt.NewOptional[Token]()
	}
	return opt.ToOptional(p.tokens[p.currentIndex+offsetAmount])
}

func (p *Parser) consume() Token {
	r := p.tokens[p.currentIndex]
	p.currentIndex++
	return r
}

func (p *Parser) tryConsume(tokType TokenType, errMsg string) (Token, error) {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == tokType {
		return p.consume(), nil
	} else {
		return Token{}, errors.New(errMsg)
	}
}
func (p *Parser) mustTryConsume(tokType TokenType) opt.Optional[Token] {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == tokType {
		return opt.ToOptional(p.consume())
	} else {
		return opt.Optional[Token]{}
	}
}

func (p *Parser) error(message string) error {
	var lineInfo LineInfo
	if p.currentIndex < len(p.tokens) {
		lineInfo = p.tokens[p.currentIndex].lineInfo
	} else {
		lineInfo = p.tokens[len(p.tokens)-1].lineInfo
	}
	return lineInfo.PositionedError(message)
}

type TypedIdentifier struct {
	ident Token
	_type NodeType
}

type NodeType interface {
	IsNodeType()
}

type NodePureType struct {
	baseType NodeBaseType
}

func (NodePureType) IsNodeType() {}

type NodePointerType struct {
	subType NodeType
}

func (NodePointerType) IsNodeType() {}

type NodeBaseType interface {
	NodeType
	IsNodeBaseType()
}

type NodeBoolType struct {
	_bool Token
}

func (NodeBoolType) IsNodeBaseType() {}
func (NodeBoolType) IsNodeType()     {}

type NodeIntType struct {
	_int Token
}

func (NodeIntType) IsNodeBaseType() {}
func (NodeIntType) IsNodeType()     {}

type NodeCharType struct {
	_char Token
}

func (NodeCharType) IsNodeBaseType() {}
func (NodeCharType) IsNodeType()     {}

type NodeProg struct {
	stmts []NodeStmt
}

type NodeStmt interface {
	IsNodeStmt()
}

type NodeStmtVarDeclare struct {
	value TypedIdentifier
}

func (NodeStmtVarDeclare) IsNodeStmt() {}

type NodeStmtVarAssign struct {
	ident Token
	expr  NodeExpr
}

func (NodeStmtVarAssign) IsNodeStmt() {}

type NodeStmtPointerAssign struct {
	ident Token
	expr  NodeExpr
}

func (NodeStmtPointerAssign) IsNodeStmt() {}

type NodeStmtIf struct {
	expr       NodeExpr
	scope      NodeScope
	elseBranch opt.Optional[NodeElse]
}

func (NodeStmtIf) IsNodeStmt() {}

type NodeStmtWhile struct {
	expr  NodeExpr
	scope NodeScope
}

func (NodeStmtWhile) IsNodeStmt() {}

type NodeStmtBreak struct {
	_break Token
}

func (NodeStmtBreak) IsNodeStmt() {}

type NodeStmtContinue struct {
	_continue Token
}

func (NodeStmtContinue) IsNodeStmt() {}

type NodeStmtFunctionDefinition struct {
	ident   Token
	params  []TypedIdentifier
	returns string
	body    NodeScope
}

func (NodeStmtFunctionDefinition) IsNodeStmt() {}

type NodeStmtReturn struct {
	returns []NodeExpr
	_return Token
}

func (NodeStmtReturn) IsNodeStmt() {}

type NodeStmtSyscall struct {
	arguments []NodeExpr
	syscall   Token
}

func (NodeStmtSyscall) IsNodeStmt() {}

type NodeExpr interface {
	IsNodeExpr()
}

type NodeIntBinExpr interface {
	NodeExpr
	IsNodeBinExpr()
}

type NodeIntBinExprAdd struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeIntBinExprAdd) IsNodeBinExpr() {}
func (NodeIntBinExprAdd) IsNodeExpr()    {}

type NodeIntBinExprSubtract struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeIntBinExprSubtract) IsNodeBinExpr() {}
func (NodeIntBinExprSubtract) IsNodeExpr()    {}

type NodeIntBinExprMultiply struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeIntBinExprMultiply) IsNodeBinExpr() {}
func (NodeIntBinExprMultiply) IsNodeExpr()    {}

type NodeIntBinExprDivide struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeIntBinExprDivide) IsNodeBinExpr() {}
func (NodeIntBinExprDivide) IsNodeExpr()    {}

type NodeIntBinExprModulo struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeIntBinExprModulo) IsNodeBinExpr() {}
func (NodeIntBinExprModulo) IsNodeExpr()    {}

type NodeIntTerm interface {
	NodeExpr
	IsNodeTerm()
}

type NodeIntTermNegativeTerm struct {
	term NodeIntTerm
}

func (NodeIntTermNegativeTerm) IsNodeTerm() {}
func (NodeIntTermNegativeTerm) IsNodeExpr() {}

type NodeIntTermLiteral struct {
	intLiteral Token
}

func (NodeIntTermLiteral) IsNodeTerm() {}
func (NodeIntTermLiteral) IsNodeExpr() {}

type NodeIntTermIdentifier struct {
	identifier Token
}

func (NodeIntTermIdentifier) IsNodeTerm() {}
func (NodeIntTermIdentifier) IsNodeExpr() {}

type NodeIntFunctionCall struct {
	ident  Token
	params []NodeExpr
}

func (NodeIntFunctionCall) IsNodeStmt() {}
func (NodeIntFunctionCall) IsNodeTerm() {}
func (NodeIntFunctionCall) IsNodeExpr() {}

type NodeIntTermRoundBracketExpr struct {
	expr NodeExpr
}

func (NodeIntTermRoundBracketExpr) IsNodeTerm() {}
func (NodeIntTermRoundBracketExpr) IsNodeExpr() {}

type NodeIntTermPointer struct {
	identifier Token
}

func (NodeIntTermPointer) IsNodeTerm() {}
func (NodeIntTermPointer) IsNodeExpr() {}

type NodeIntTermPointerDereference struct {
	identifier Token
}

func (NodeIntTermPointerDereference) IsNodeTerm() {}
func (NodeIntTermPointerDereference) IsNodeExpr() {}

type NodeScope struct {
	stmts []NodeStmt
}

func (NodeScope) IsNodeStmt() {}

type NodeElse interface {
	IsNodeElif()
}

type NodeElseElif struct {
	ifStmt NodeStmtIf
}

func (NodeElseElif) IsNodeElif() {}

type NodeElseScope struct {
	scope NodeScope
}

func (NodeElseScope) IsNodeElif() {}
