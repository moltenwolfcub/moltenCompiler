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

func (p *Parser) ParseProg() (opt.Optional[NodeProg], error) {
	node := NodeProg{
		[]NodeStmt{},
	}
	for p.peek().HasValue() {
		stmt, err := p.ParseStmt()
		if err != nil {
			return opt.Optional[NodeProg]{}, err
		}

		if stmt.HasValue() {
			node.stmts = append(node.stmts, stmt.MustGetValue())
		} else {
			return opt.Optional[NodeProg]{}, p.error("invalid statment")
		}
	}
	return opt.ToOptional(node), nil
}

func (p *Parser) ParseStmt() (opt.Optional[NodeStmt], error) {
	if p.mustTryConsume(exit).HasValue() {
		_, err := p.tryConsume(openRoundBracket, "expected '(' after 'exit'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		var node NodeStmtExit

		nodeExpr, err := p.ParseExpr()
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		if nodeExpr.HasValue() {
			node = NodeStmtExit{nodeExpr.MustGetValue()}

		} else if p.mustTryConsume(closeRoundBracket).HasValue() {
			node = NodeStmtExit{NodeExpr{NodeTerm{NodeTermIntLiteral{Token{
				tokenType: intLiteral,
				value:     opt.ToOptional("0"),
			}}}}}
		} else {
			return opt.Optional[NodeStmt]{}, p.error("invalid expression for exit")
		}

		_, err = p.tryConsume(closeRoundBracket, "missing ')'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		_, err = p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		return opt.ToOptional(NodeStmt{node}), nil
	} else if p.mustTryConsume(_var).HasValue() {

		tok, err := p.tryConsume(identifier, "expected variable name after `var`")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		node := NodeStmtVarDeclare{tok}

		_, err = p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		return opt.ToOptional(NodeStmt{node}), nil
	} else if tok := p.mustTryConsume(identifier); tok.HasValue() {
		node := NodeStmtVarAssign{
			ident: tok.MustGetValue(),
		}
		_, err := p.tryConsume(equals, "expected '=' after variable name for assignment")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		expr, err := p.ParseExpr()
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		if expr.HasValue() {
			node.expr = expr.MustGetValue()
		} else {
			return opt.Optional[NodeStmt]{}, p.error("invalid expression for variable assigment")
		}

		_, err = p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		return opt.ToOptional(NodeStmt{node}), nil
	} else if p.peek().MustGetValue().tokenType == openCurlyBracket {
		scope, err := p.ParseScope()
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		if !scope.HasValue() {
			panic(errors.New("invalid scope"))
		}
		return opt.ToOptional(NodeStmt{scope.MustGetValue()}), nil

	} else if p.peek().MustGetValue().tokenType == _if {
		ifStmt, err := p.ParseIf()
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		if !ifStmt.HasValue() {
			panic(errors.New("invalid if statement"))
		}
		return opt.ToOptional(NodeStmt{ifStmt.MustGetValue()}), nil

	} else if tok := p.mustTryConsume(while); tok.HasValue() {
		node := NodeStmtWhile{}

		_, err := p.tryConsume(openRoundBracket, "Expected '('")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		expr, err := p.ParseExpr()
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		if expr.HasValue() {
			node.expr = expr.MustGetValue()
		} else {
			panic(errors.New("invalid expression"))
		}

		_, err = p.tryConsume(closeRoundBracket, "Expected ')'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}

		scope, err := p.ParseScope()
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		if scope.HasValue() {
			node.scope = scope.MustGetValue()
		} else {
			panic(errors.New("invalid if statement, expected scope"))
		}
		return opt.ToOptional(NodeStmt{node}), nil

	} else if tok := p.mustTryConsume(_break); tok.HasValue() {
		_, err := p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		return opt.ToOptional(NodeStmt{NodeStmtBreak{}}), nil

	} else if tok := p.mustTryConsume(_continue); tok.HasValue() {
		_, err := p.tryConsume(semiColon, "missing ';'")
		if err != nil {
			return opt.Optional[NodeStmt]{}, err
		}
		return opt.ToOptional(NodeStmt{NodeStmtContinue{}}), nil

	} else {
		return opt.Optional[NodeStmt]{}, nil
	}
}

func (p *Parser) ParseTerm() (opt.Optional[NodeTerm], error) {
	if tok := p.mustTryConsume(intLiteral); tok.HasValue() {
		return opt.ToOptional(NodeTerm{NodeTermIntLiteral{tok.MustGetValue()}}), nil
	} else if tok := p.mustTryConsume(identifier); tok.HasValue() {
		return opt.ToOptional(NodeTerm{NodeTermIdentifier{tok.MustGetValue()}}), nil
	} else if p.mustTryConsume(openRoundBracket).HasValue() {
		expr, err := p.ParseExpr()
		if err != nil {
			return opt.Optional[NodeTerm]{}, err
		}
		if !expr.HasValue() {
			panic(errors.New("expected expression"))
		}
		_, err = p.tryConsume(closeRoundBracket, "expected ')'")
		if err != nil {
			return opt.Optional[NodeTerm]{}, err
		}
		return opt.ToOptional(NodeTerm{NodeTermRoundBracketExpr{expr.MustGetValue()}}), nil
	}
	return opt.Optional[NodeTerm]{}, nil
}

func (p *Parser) ParseScope() (opt.Optional[NodeScope], error) {
	if !p.mustTryConsume(openCurlyBracket).HasValue() {
		return opt.Optional[NodeScope]{}, nil
	}

	var scope NodeScope
	for {
		stmt, err := p.ParseStmt()
		if err != nil {
			return opt.Optional[NodeScope]{}, err
		}
		if !stmt.HasValue() {
			break
		}

		scope.stmts = append(scope.stmts, stmt.MustGetValue())
	}
	_, err := p.tryConsume(closeCurlyBracket, "expected '}'")
	if err != nil {
		return opt.Optional[NodeScope]{}, err
	}

	return opt.ToOptional(scope), nil
}

func (p *Parser) ParseIf() (opt.Optional[NodeStmtIf], error) {
	if !p.mustTryConsume(_if).HasValue() {
		return opt.Optional[NodeStmtIf]{}, nil
	}

	node := NodeStmtIf{}

	_, err := p.tryConsume(openRoundBracket, "Expected '('")
	if err != nil {
		return opt.Optional[NodeStmtIf]{}, err
	}

	expr, err := p.ParseExpr()
	if err != nil {
		return opt.Optional[NodeStmtIf]{}, err
	}
	if expr.HasValue() {
		node.expr = expr.MustGetValue()
	} else {
		panic(errors.New("invalid expression"))
	}
	_, err = p.tryConsume(closeRoundBracket, "Expected ')'")
	if err != nil {
		return opt.Optional[NodeStmtIf]{}, err
	}

	scope, err := p.ParseScope()
	if err != nil {
		return opt.Optional[NodeStmtIf]{}, err
	}
	if scope.HasValue() {
		node.scope = scope.MustGetValue()
	} else {
		panic(errors.New("invalid if statement, expected scope"))
	}

	if p.mustTryConsume(_else).HasValue() {
		node.elseBranch, err = p.ParseElse()
		if err != nil {
			return opt.Optional[NodeStmtIf]{}, err
		}
	}

	return opt.ToOptional(node), nil
}

func (p *Parser) ParseElse() (opt.Optional[NodeElse], error) {
	node := NodeElse{}

	ifStmt, err := p.ParseIf()
	if err != nil {
		return opt.Optional[NodeElse]{}, err
	}
	if ifStmt.HasValue() {
		elifNode := NodeElseElif{ifStmt.MustGetValue()}
		node.variant = elifNode
		return opt.ToOptional(node), nil
	} else {
		scope, err := p.ParseScope()
		if err != nil {
			return opt.Optional[NodeElse]{}, err
		}
		if scope.HasValue() {
			scopeNode := NodeElseScope{scope.MustGetValue()}
			node.variant = scopeNode
			return opt.ToOptional(node), nil
		}
	}
	return opt.Optional[NodeElse]{}, nil
}

// based off of this principle and algorithm:
// https://eli.thegreenplace.net/2012/08/02/parsing-expressions-by-precedence-climbing
func (p *Parser) ParseExpr(minPrecedence ...int) (opt.Optional[NodeExpr], error) {
	var minPrec int
	if len(minPrecedence) == 1 {
		minPrec = minPrecedence[0]
	} else {
		minPrec = 0
	}

	lhsTerm, err := p.ParseTerm()
	if err != nil {
		return opt.Optional[NodeExpr]{}, err
	}

	if !lhsTerm.HasValue() {
		return opt.Optional[NodeExpr]{}, nil
	}

	lhsExpr := NodeExpr{lhsTerm.MustGetValue()}

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

		rhsExpr, err := p.ParseExpr(nextMinPrec)
		if err != nil {
			return opt.Optional[NodeExpr]{}, err
		}
		if !rhsExpr.HasValue() {
			panic(errors.New("unable to parse expression"))
		}

		expr := NodeBinExpr{}
		switch op.tokenType {
		case plus:
			add := NodeBinExprAdd{
				left:  lhsExpr,
				right: rhsExpr.MustGetValue(),
			}
			expr.variant = add
		case asterisk:
			multiply := NodeBinExprMultiply{
				left:  lhsExpr,
				right: rhsExpr.MustGetValue(),
			}
			expr.variant = multiply
		case minus:
			subtract := NodeBinExprSubtract{
				left:  lhsExpr,
				right: rhsExpr.MustGetValue(),
			}
			expr.variant = subtract
		case fslash:
			divide := NodeBinExprDivide{
				left:  lhsExpr,
				right: rhsExpr.MustGetValue(),
			}
			expr.variant = divide
		}
		lhsExpr.variant = expr

	}
	return opt.ToOptional(lhsExpr), nil
}

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
		return Token{}, p.error(errMsg)
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
	currentToken := p.tokens[p.currentIndex]

	return fmt.Errorf("%s:%d:%d: %s", currentToken.file, currentToken.line, currentToken.col, message)
}

type NodeProg struct {
	stmts []NodeStmt
}

type NodeStmt struct {
	variant interface {
		IsNodeStmt()
	}
}

type NodeStmtExit struct {
	expr NodeExpr
}

func (NodeStmtExit) IsNodeStmt() {}

type NodeStmtVarDeclare struct {
	ident Token
}

func (NodeStmtVarDeclare) IsNodeStmt() {}

type NodeStmtVarAssign struct {
	ident Token
	expr  NodeExpr
}

func (NodeStmtVarAssign) IsNodeStmt() {}

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

type NodeStmtBreak struct{}

func (NodeStmtBreak) IsNodeStmt() {}

type NodeStmtContinue struct{}

func (NodeStmtContinue) IsNodeStmt() {}

type NodeExpr struct {
	variant interface {
		IsNodeExpr()
	}
}

type NodeBinExpr struct {
	variant interface {
		IsNodeBinExpr()
	}
}

func (NodeBinExpr) IsNodeExpr() {}

type NodeBinExprAdd struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeBinExprAdd) IsNodeBinExpr() {}

type NodeBinExprSubtract struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeBinExprSubtract) IsNodeBinExpr() {}

type NodeBinExprMultiply struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeBinExprMultiply) IsNodeBinExpr() {}

type NodeBinExprDivide struct {
	left  NodeExpr
	right NodeExpr
}

func (NodeBinExprDivide) IsNodeBinExpr() {}

type NodeTerm struct {
	variant interface {
		IsNodeTerm()
	}
}

func (NodeTerm) IsNodeExpr() {}

type NodeTermIntLiteral struct {
	intLiteral Token
}

func (n NodeTermIntLiteral) IsNodeTerm() {}

type NodeTermIdentifier struct {
	identifier Token
}

func (NodeTermIdentifier) IsNodeTerm() {}

type NodeTermRoundBracketExpr struct {
	expr NodeExpr
}

func (NodeTermRoundBracketExpr) IsNodeTerm() {}

type NodeScope struct {
	stmts []NodeStmt
}

func (NodeScope) IsNodeStmt() {}

type NodeElse struct {
	variant interface {
		IsNodeElif()
	}
}

type NodeElseElif struct {
	ifStmt NodeStmtIf
}

func (NodeElseElif) IsNodeElif() {}

type NodeElseScope struct {
	scope NodeScope
}

func (NodeElseScope) IsNodeElif() {}
