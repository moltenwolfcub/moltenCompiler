package main

import (
	"errors"
	"fmt"

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

			expr, err := p.ParseExpr()
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

		expr, err := p.ParseExpr()
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

		node.expr, err = p.ParseIntExpr() //TODO: change to bool expr
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
			expr, err := p.ParseExpr()
			if err == errMissingIntExpr {
				break
			} else if err != nil {
				return NodeStmtReturn{}, err
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
			expr, err := p.ParseExpr()
			if err == errMissingIntExpr {
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

func (p *Parser) ParseExpr() (NodeExpr, error) {
	boolExpr, err := p.ParseBoolExpr()
	if err == nil {
		return boolExpr, nil
	}

	intExpr, err := p.ParseIntExpr()
	if err == nil {
		return intExpr, nil
	}

	return nil, errMissingExpr
}

var errMissingExpr error = errors.New("expected expression")

// based off of this principle and algorithm:
// https://eli.thegreenplace.net/2012/08/02/parsing-expressions-by-precedence-climbing
func (p *Parser) ParseIntExpr(minPrecedence ...int) (NodeIntExpr, error) {
	var minPrec int
	if len(minPrecedence) == 1 {
		minPrec = minPrecedence[0]
	} else {
		minPrec = 0
	}

	lhsTerm, err := p.ParseIntTerm()
	if err == errMissingIntTerm {
		return nil, errMissingIntExpr
	}
	if err != nil {
		return nil, err
	}

	lhsExpr := NodeIntExpr(lhsTerm)

	for {
		currentToken := p.peek()
		if !currentToken.HasValue() {
			break
		}
		currentPrec := currentToken.MustGetValue().tokenType.GetIntBinPrec()
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

var errMissingIntExpr error = errors.New("expected integer expression")

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
		return NodeTermPointerDereference{variable}, nil
	}
	return nil, errMissingIntTerm
}

var errMissingIntTerm error = errors.New("expected integer term but couldn't find one")

func (p *Parser) ParseBoolExpr(minPrecedence ...int) (NodeBoolExpr, error) {
	var minPrec int
	if len(minPrecedence) == 1 {
		minPrec = minPrecedence[0]
	} else {
		minPrec = 0
	}

	lhsTerm, err := p.ParseBoolTerm()
	if err == errMissingBoolTerm {
		return nil, errMissingBoolExpr
	}
	if err != nil {
		return nil, err
	}

	lhsExpr := NodeBoolExpr(lhsTerm)

	for {
		currentToken := p.peek()
		if !currentToken.HasValue() {
			break
		}
		currentPrec := currentToken.MustGetValue().tokenType.GetBoolBinPrec()
		if !currentPrec.HasValue() || currentPrec.MustGetValue() < minPrec { //prolly meant to be <=
			break
		}
		op := p.consume()
		if !p.peek().HasValue() || p.peek().MustGetValue().tokenType != op.tokenType {
			return nil, errors.New("boolean expressions must have a double operator E.G.(&&, ||)")
		}
		p.consume()

		nextMinPrec := currentPrec.MustGetValue() + 1

		rhsExpr, err := p.ParseBoolExpr(nextMinPrec)
		if err != nil {
			return nil, err
		}

		var expr NodeBoolBinExpr
		switch op.tokenType {
		case ampersand:
			expr = NodeBoolBinExprAnd{
				left:  lhsExpr,
				right: rhsExpr,
			}
		case pipe:
			expr = NodeBoolBinExprOr{
				left:  lhsExpr,
				right: rhsExpr,
			}
		}
		lhsExpr = expr

	}
	return lhsExpr, nil
}

var errMissingBoolExpr error = errors.New("expected boolean expression")

func (p *Parser) ParseBoolTerm() (NodeBoolTerm, error) {
	BoolTermWrapper := func() (NodeBoolTerm, error) {
		if p.mustTryConsume(exclamation).HasValue() {
			term, err := p.ParseBoolTerm()
			if err != nil {
				return nil, err
			}
			return NodeBoolTermNotTerm{term}, nil
		} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == identifier {
			if p.peek(1).MustGetValue().tokenType == openRoundBracket {
				return p.ParseFuncCall()
			} else {
				return NodeBoolTermIdentifier{p.consume()}, nil
			}
		} else if p.mustTryConsume(openRoundBracket).HasValue() {
			expr, err := p.ParseBoolExpr()
			if err != nil {
				return nil, err
			}
			_, err = p.tryConsume(closeRoundBracket, "expected ')'")
			if err != nil {
				return nil, err
			}
			return NodeBoolTermRoundBracketExpr{expr}, nil
		} else if p.mustTryConsume(asterisk).HasValue() {
			variable, err := p.tryConsume(identifier, "expected variable identifier after '*'")
			if err != nil {
				return nil, err
			}
			return NodeTermPointerDereference{variable}, nil
		}
		return nil, errMissingBoolTerm
	}

	parsedBool, err := BoolTermWrapper()
	if err == nil {
		relativeOp, err := p.ParseRelativeOperator()
		if err != nil {
			return parsedBool, nil //if no relOp
		}

		rhsParsedBool, err := BoolTermWrapper()
		if err != nil {
			return nil, errors.New("found a relative operator but no term following it") // err if relOp but no rhs
		}
		return NodeBoolComparisonBool{
			left:  parsedBool,
			right: rhsParsedBool,
			op:    relativeOp,
		}, nil
	} else if err != errMissingBoolTerm {
		return nil, err
	}

	parsedInt, err := p.ParseIntExpr()
	if err == nil {
		relativeOp, err := p.ParseRelativeOperator()
		if err != nil {
			return nil, errors.New("found intTerm but no relative operator. can't use intTerm in boolean expression")
		}

		t := p.peek()
		fmt.Println(t)

		rhsParsedInt, err := p.ParseIntExpr()
		if err != nil {
			return nil, errors.New("found a relative operator but no term following it")
		}
		return NodeBoolComparisonInt{
			left:  parsedInt,
			right: rhsParsedInt,
			op:    relativeOp,
		}, nil
	} else if err != errMissingIntExpr {
		return nil, err
	}

	return nil, errMissingBoolTerm
	//TODO: return p.currentIndex to initial value in event of failure
}

var errMissingBoolTerm error = errors.New("expected boolean term but couldn't find one")

func (p *Parser) ParseRelativeOperator() (NodeRelativeOperator, error) {
	if tok := p.mustTryConsume(equals); tok.HasValue() {
		if p.mustTryConsume(equals).HasValue() {
			return NodeRelativeOpEqual{equal: tok.MustGetValue()}, nil
		}
	} else if tok := p.mustTryConsume(exclamation); tok.HasValue() {
		if p.mustTryConsume(equals).HasValue() {
			return NodeRelativeOpNotEqual{notEqual: tok.MustGetValue()}, nil
		}
	} else if tok := p.mustTryConsume(openTriangleBracket); tok.HasValue() {
		if p.mustTryConsume(equals).HasValue() {
			return NodeRelativeOpLessThanOrEqual{lessThanOrEqual: tok.MustGetValue()}, nil
		}
		return NodeRelativeOpLessThan{lessThan: tok.MustGetValue()}, nil
	} else if tok := p.mustTryConsume(closeTriangleBracket); tok.HasValue() {
		if p.mustTryConsume(equals).HasValue() {
			return NodeRelativeOpGreaterThanOrEqual{greaterThanOrEqual: p.consume()}, nil
		}
		return NodeRelativeOpGreaterThan{greaterThan: tok.MustGetValue()}, nil
	}

	return nil, errMissingRelativeOperator
}

var errMissingRelativeOperator error = errors.New("expected relative operator but couldn't find one")

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

	node.expr, err = p.ParseIntExpr() //TODO: change to bool expr
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

func (p *Parser) ParseFuncCall() (NodeFunctionCall, error) {

	node := NodeFunctionCall{
		ident: p.consume(),
	}

	_, err := p.tryConsume(openRoundBracket, "missing '('")
	if err != nil {
		return NodeFunctionCall{}, err
	}

	for {
		expr, err := p.ParseExpr()
		if err == errMissingIntExpr {
			break
		} else if err != nil {
			return NodeFunctionCall{}, err
		}
		node.params = append(node.params, expr)

		_, err = p.tryConsume(comma, "optional so this should never error")
		if err != nil {
			break
		}
	}

	_, err = p.tryConsume(closeRoundBracket, "missing ')'")
	if err != nil {
		return NodeFunctionCall{}, err
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

// region types
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

// endregion

type NodeProg struct {
	stmts []NodeStmt
}

// region statements

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

//endregion

//region expressions

type NodeExpr interface {
	IsNodeExpr()
}

type NodeFunctionCall struct {
	ident  Token
	params []NodeExpr
}

func (NodeFunctionCall) IsNodeStmt()     {}
func (NodeFunctionCall) IsNodeIntTerm()  {}
func (NodeFunctionCall) IsNodeIntExpr()  {}
func (NodeFunctionCall) IsNodeBoolTerm() {}
func (NodeFunctionCall) IsNodeBoolExpr() {}
func (NodeFunctionCall) IsNodeExpr()     {}

type NodeTermPointerDereference struct {
	identifier Token
}

func (NodeTermPointerDereference) IsNodeIntTerm()  {}
func (NodeTermPointerDereference) IsNodeIntExpr()  {}
func (NodeTermPointerDereference) IsNodeBoolTerm() {}
func (NodeTermPointerDereference) IsNodeBoolExpr() {}
func (NodeTermPointerDereference) IsNodeExpr()     {}

//region intExprs

type NodeIntExpr interface {
	NodeExpr
	IsNodeIntExpr()
}

type NodeIntBinExpr interface {
	NodeIntExpr
	IsNodeIntBinExpr()
}

type NodeIntBinExprAdd struct {
	left  NodeIntExpr
	right NodeIntExpr
}

func (NodeIntBinExprAdd) IsNodeIntBinExpr() {}
func (NodeIntBinExprAdd) IsNodeIntExpr()    {}
func (NodeIntBinExprAdd) IsNodeExpr()       {}

type NodeIntBinExprSubtract struct {
	left  NodeIntExpr
	right NodeIntExpr
}

func (NodeIntBinExprSubtract) IsNodeIntBinExpr() {}
func (NodeIntBinExprSubtract) IsNodeIntExpr()    {}
func (NodeIntBinExprSubtract) IsNodeExpr()       {}

type NodeIntBinExprMultiply struct {
	left  NodeIntExpr
	right NodeIntExpr
}

func (NodeIntBinExprMultiply) IsNodeIntBinExpr() {}
func (NodeIntBinExprMultiply) IsNodeIntExpr()    {}
func (NodeIntBinExprMultiply) IsNodeExpr()       {}

type NodeIntBinExprDivide struct {
	left  NodeIntExpr
	right NodeIntExpr
}

func (NodeIntBinExprDivide) IsNodeIntBinExpr() {}
func (NodeIntBinExprDivide) IsNodeIntExpr()    {}
func (NodeIntBinExprDivide) IsNodeExpr()       {}

type NodeIntBinExprModulo struct {
	left  NodeIntExpr
	right NodeIntExpr
}

func (NodeIntBinExprModulo) IsNodeIntBinExpr() {}
func (NodeIntBinExprModulo) IsNodeIntExpr()    {}
func (NodeIntBinExprModulo) IsNodeExpr()       {}

type NodeIntTerm interface {
	NodeIntExpr
	IsNodeIntTerm()
}

type NodeIntTermNegativeTerm struct {
	term NodeIntTerm
}

func (NodeIntTermNegativeTerm) IsNodeIntTerm() {}
func (NodeIntTermNegativeTerm) IsNodeIntExpr() {}
func (NodeIntTermNegativeTerm) IsNodeExpr()    {}

type NodeIntTermLiteral struct {
	intLiteral Token
}

func (NodeIntTermLiteral) IsNodeIntTerm() {}
func (NodeIntTermLiteral) IsNodeIntExpr() {}
func (NodeIntTermLiteral) IsNodeExpr()    {}

type NodeIntTermIdentifier struct {
	identifier Token
}

func (NodeIntTermIdentifier) IsNodeIntTerm() {}
func (NodeIntTermIdentifier) IsNodeIntExpr() {}
func (NodeIntTermIdentifier) IsNodeExpr()    {}

type NodeIntTermRoundBracketExpr struct {
	expr NodeIntExpr
}

func (NodeIntTermRoundBracketExpr) IsNodeIntTerm() {}
func (NodeIntTermRoundBracketExpr) IsNodeIntExpr() {}
func (NodeIntTermRoundBracketExpr) IsNodeExpr()    {}

type NodeIntTermPointer struct {
	identifier Token
}

func (NodeIntTermPointer) IsNodeIntTerm() {}
func (NodeIntTermPointer) IsNodeIntExpr() {}
func (NodeIntTermPointer) IsNodeExpr()    {}

//endregion

//region boolExprs

type NodeBoolExpr interface {
	NodeExpr
	IsNodeBoolExpr()
}

type NodeBoolBinExpr interface {
	NodeBoolExpr
	IsNodeBoolBinExpr()
}

type NodeBoolBinExprAnd struct {
	left  NodeBoolExpr
	right NodeBoolExpr
}

func (NodeBoolBinExprAnd) IsNodeBoolBinExpr() {}
func (NodeBoolBinExprAnd) IsNodeBoolExpr()    {}
func (NodeBoolBinExprAnd) IsNodeExpr()        {}

type NodeBoolBinExprOr struct {
	left  NodeBoolExpr
	right NodeBoolExpr
}

func (NodeBoolBinExprOr) IsNodeBoolBinExpr() {}
func (NodeBoolBinExprOr) IsNodeBoolExpr()    {}
func (NodeBoolBinExprOr) IsNodeExpr()        {}

type NodeBoolTerm interface {
	NodeBoolExpr
	IsNodeBoolTerm()
}

type NodeBoolTermNotTerm struct {
	term NodeBoolTerm
}

func (NodeBoolTermNotTerm) IsNodeBoolTerm() {}
func (NodeBoolTermNotTerm) IsNodeBoolExpr() {}
func (NodeBoolTermNotTerm) IsNodeExpr()     {}

type NodeBoolTermIdentifier struct {
	identifier Token
}

func (NodeBoolTermIdentifier) IsNodeBoolTerm() {}
func (NodeBoolTermIdentifier) IsNodeBoolExpr() {}
func (NodeBoolTermIdentifier) IsNodeExpr()     {}

type NodeBoolTermRoundBracketExpr struct {
	expr NodeBoolExpr
}

func (NodeBoolTermRoundBracketExpr) IsNodeBoolTerm() {}
func (NodeBoolTermRoundBracketExpr) IsNodeBoolExpr() {}
func (NodeBoolTermRoundBracketExpr) IsNodeExpr()     {}

type NodeBoolComparisonBool struct {
	left  NodeBoolTerm
	right NodeBoolTerm
	op    NodeRelativeOperator
}

func (NodeBoolComparisonBool) IsNodeBoolTerm() {}
func (NodeBoolComparisonBool) IsNodeBoolExpr() {}
func (NodeBoolComparisonBool) IsNodeExpr()     {}

type NodeBoolComparisonInt struct {
	left  NodeIntExpr
	right NodeIntExpr
	op    NodeRelativeOperator
}

func (NodeBoolComparisonInt) IsNodeBoolTerm() {}
func (NodeBoolComparisonInt) IsNodeBoolExpr() {}
func (NodeBoolComparisonInt) IsNodeExpr()     {}

// endregion

// endregion

type NodeRelativeOperator interface {
	IsRelativeOperator()
}

type NodeRelativeOpEqual struct {
	equal Token
}

func (NodeRelativeOpEqual) IsRelativeOperator() {}

type NodeRelativeOpNotEqual struct {
	notEqual Token
}

func (NodeRelativeOpNotEqual) IsRelativeOperator() {}

type NodeRelativeOpLessThan struct {
	lessThan Token
}

func (NodeRelativeOpLessThan) IsRelativeOperator() {}

type NodeRelativeOpGreaterThan struct {
	greaterThan Token
}

func (NodeRelativeOpGreaterThan) IsRelativeOperator() {}

type NodeRelativeOpLessThanOrEqual struct {
	lessThanOrEqual Token
}

func (NodeRelativeOpLessThanOrEqual) IsRelativeOperator() {}

type NodeRelativeOpGreaterThanOrEqual struct {
	greaterThanOrEqual Token
}

func (NodeRelativeOpGreaterThanOrEqual) IsRelativeOperator() {}

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
