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
	if p.peek().MustGetValue().tokenType == exit {
		if !p.peek(1).HasValue() || p.peek(1).MustGetValue().tokenType != openRoundBracket {
			panic(errors.New("expected '(' after 'exit'"))
		}

		p.consume()
		p.consume()

		var node NodeStmtExit

		if nodeExpr := p.ParseExpr(); nodeExpr.HasValue() {
			node = NodeStmtExit{expr: nodeExpr.MustGetValue()}

		} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == closeRoundBracket {
			node = NodeStmtExit{
				expr: NodeExpr{NodeExprIntLiteral{intLiteral: Token{
					tokenType: intLiteral,
					value:     opt.ToOptional("0"),
				}}},
			}
		} else {
			panic(errors.New("invalid expression"))
		}

		if p.peek().HasValue() && p.peek().MustGetValue().tokenType == closeRoundBracket {
			p.consume()
		} else {
			panic(errors.New("missing ')'"))
		}

		if p.peek().HasValue() && p.peek().MustGetValue().tokenType == semiColon {
			p.consume()
		} else {
			panic(errors.New("missing ';'"))
		}

		return opt.ToOptional(NodeStmt{variant: node})
	} else if p.peek().MustGetValue().tokenType == _var {
		p.consume()

		var node NodeStmtVarDeclare

		if p.peek().HasValue() && p.peek().MustGetValue().tokenType == identifier {
			node = NodeStmtVarDeclare{ident: p.consume()}
		} else {
			panic(errors.New("expected variable name after `var`"))
		}

		if p.peek().HasValue() && p.peek().MustGetValue().tokenType == semiColon {
			p.consume()
		} else {
			panic(errors.New("missing ';'"))
		}

		return opt.ToOptional(NodeStmt{variant: node})
	} else if p.peek().MustGetValue().tokenType == identifier &&
		p.peek(1).HasValue() && p.peek(1).MustGetValue().tokenType == equals {

		node := NodeStmtVarAssign{
			ident: p.consume(),
		}
		p.consume()

		if expr := p.ParseExpr(); expr.HasValue() {
			node.expr = expr.MustGetValue()
		} else {
			panic(errors.New("invalid expression"))
		}

		if p.peek().HasValue() && p.peek().MustGetValue().tokenType == semiColon {
			p.consume()
		} else {
			panic(errors.New("missing ';'"))
		}

		return opt.ToOptional(NodeStmt{variant: node})
	} else {
		return opt.Optional[NodeStmt]{}
	}
}

func (p *Parser) ParseExpr() opt.Optional[NodeExpr] {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == intLiteral {
		return opt.ToOptional(NodeExpr{variant: NodeExprIntLiteral{intLiteral: p.consume()}})

	} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == identifier {
		return opt.ToOptional(NodeExpr{variant: NodeExprIdentifier{identifier: p.consume()}})
	}
	return opt.Optional[NodeExpr]{}
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

func (n NodeStmtExit) IsNodeStmt() {}

type NodeStmtVarDeclare struct {
	ident Token
}

func (n NodeStmtVarDeclare) IsNodeStmt() {}

type NodeStmtVarAssign struct {
	ident Token
	expr  NodeExpr
}

func (n NodeStmtVarAssign) IsNodeStmt() {}

type NodeExpr struct {
	variant interface {
		IsNodeExpr()
	}
}

type NodeExprIntLiteral struct {
	intLiteral Token
}

func (n NodeExprIntLiteral) IsNodeExpr() {}

type NodeExprIdentifier struct {
	identifier Token
}

func (n NodeExprIdentifier) IsNodeExpr() {}
