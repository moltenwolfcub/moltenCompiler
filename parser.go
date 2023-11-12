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
	}
	if tok := p.tryConsume(identifier); tok.HasValue() {
		return opt.ToOptional(NodeTerm{NodeTermIdentifier{tok.MustGetValue()}})
	}
	return opt.Optional[NodeTerm]{}
}

func (p *Parser) ParseExpr() opt.Optional[NodeExpr] {
	if term := p.ParseTerm(); term.HasValue() {
		if p.peek().HasValue() && p.peek().MustGetValue().tokenType == plus {
			return opt.ToOptional(NodeExpr{p.ParseBinExpr(term.MustGetValue()).MustGetValue()})
		} else {
			return opt.ToOptional(NodeExpr{term.MustGetValue()})
		}
	}
	return opt.Optional[NodeExpr]{}
}

func (p *Parser) ParseBinExpr(lhs NodeTerm) opt.Optional[NodeBinExpr] {
	if p.tryConsume(plus).HasValue() {
		node := NodeBinExprAdd{}
		node.left = NodeExpr{lhs}
		if rhs := p.ParseExpr(); rhs.HasValue() {
			node.right = rhs.MustGetValue()
		} else {
			panic(errors.New("expected expression"))
		}
		return opt.ToOptional(NodeBinExpr{node})
	} else {
		panic(errors.New("unsupported binary operator"))
	}
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
