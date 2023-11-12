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

func (p *Parser) ParseProg() opt.Optional[NodeProg] {
	node := NodeProg{
		[]NodeStmt{},
	}
	for p.peek().HasValue() {
		if stmt := p.ParseStmt(); stmt.HasValue() {
			node.stmts = append(node.stmts, stmt.MustGetValue())
		} else {
			panic(errors.New("invalid statment"))
		}
	}
	return opt.ToOptional(node)
}

func (p *Parser) ParseStmt() opt.Optional[NodeStmt] {
	if p.tryConsume(exit).HasValue() {
		p.mustTryConsume(openRoundBracket, "expected '(' after 'exit'")

		var node NodeStmtExit

		if nodeExpr := p.ParseExpr(); nodeExpr.HasValue() {
			node = NodeStmtExit{nodeExpr.MustGetValue()}

		} else if p.tryConsume(closeRoundBracket).HasValue() {
			node = NodeStmtExit{NodeExpr{NodeTerm{NodeTermIntLiteral{Token{
				tokenType: intLiteral,
				value:     opt.ToOptional("0"),
			}}}}}
		} else {
			panic(errors.New("invalid expression"))
		}

		p.mustTryConsume(closeRoundBracket, "missing ')'")
		p.mustTryConsume(semiColon, "missing ';'")

		return opt.ToOptional(NodeStmt{node})
	} else if p.tryConsume(_var).HasValue() {

		tok := p.mustTryConsume(identifier, "expected variable name after `var`")
		node := NodeStmtVarDeclare{tok}

		p.mustTryConsume(semiColon, "missing ';'")

		return opt.ToOptional(NodeStmt{node})
	} else if tok := p.tryConsume(identifier); tok.HasValue() {
		node := NodeStmtVarAssign{
			ident: tok.MustGetValue(),
		}
		p.mustTryConsume(equals, "expected '=' after variable name for assignment")

		if expr := p.ParseExpr(); expr.HasValue() {
			node.expr = expr.MustGetValue()
		} else {
			panic(errors.New("invalid expression"))
		}

		p.mustTryConsume(semiColon, "missing ';'")

		return opt.ToOptional(NodeStmt{node})
	} else {
		return opt.Optional[NodeStmt]{}
	}
}

func (p *Parser) ParseTerm() opt.Optional[NodeTerm] {
	if tok := p.tryConsume(intLiteral); tok.HasValue() {
		return opt.ToOptional(NodeTerm{NodeTermIntLiteral{tok.MustGetValue()}})
	} else if tok := p.tryConsume(identifier); tok.HasValue() {
		return opt.ToOptional(NodeTerm{NodeTermIdentifier{tok.MustGetValue()}})
	} else if p.tryConsume(openRoundBracket).HasValue() {
		expr := p.ParseExpr()
		if !expr.HasValue() {
			panic(errors.New("expected expression"))
		}
		p.mustTryConsume(closeRoundBracket, "expected ')'")
		return opt.ToOptional(NodeTerm{NodeTermRoundBracketExpr{expr.MustGetValue()}})
	}
	return opt.Optional[NodeTerm]{}
}

// based off of this principle and algorithm:
// https://eli.thegreenplace.net/2012/08/02/parsing-expressions-by-precedence-climbing
func (p *Parser) ParseExpr(minPrecedence ...int) opt.Optional[NodeExpr] {
	var minPrec int
	if len(minPrecedence) == 1 {
		minPrec = minPrecedence[0]
	} else {
		minPrec = 0
	}

	lhsTerm := p.ParseTerm()
	if !lhsTerm.HasValue() {
		return opt.Optional[NodeExpr]{}
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

		rhsExpr := p.ParseExpr(nextMinPrec)
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
	return opt.ToOptional(lhsExpr)
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

func (p *Parser) mustTryConsume(tokType TokenType, errMsg string) Token {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == tokType {
		return p.consume()
	} else {
		panic(errors.New(errMsg))
	}
}
func (p *Parser) tryConsume(tokType TokenType) opt.Optional[Token] {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == tokType {
		return opt.ToOptional(p.consume())
	} else {
		return opt.Optional[Token]{}
	}
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
