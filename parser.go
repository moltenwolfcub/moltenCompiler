package main

import (
	"errors"

	opt "github.com/moltenwolfcub/moltenCompiler/optional"
)

type NodeExit struct {
	expression NodeExpr
}
type NodeExpr struct {
	intLiteral Token
}

type Parser struct {
	tokens       []Token
	currentIndex int
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p Parser) Parse() (opt.Optional[NodeExit], error) {
	node := opt.NewOptional[NodeExit]()
	for p.peek().HasValue() {
		if p.peek().MustGetValue().tokenType == exit {
			p.consume()

			if nodeExpr := p.ParseExpr(); nodeExpr.HasValue() {
				node.SetValue(NodeExit{expression: nodeExpr.MustGetValue()})
			} else {
				return opt.NewOptional[NodeExit](), errors.New("invalid expression")
			}

			if p.peek().HasValue() && p.peek().MustGetValue().tokenType == semiColon {
				p.consume()
			} else {
				return opt.NewOptional[NodeExit](), errors.New("missing ';'")
			}
		}
	}
	return node, nil
}

func (p Parser) ParseExpr() opt.Optional[NodeExpr] {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == intLiteral {
		return opt.ToOptional(NodeExpr{intLiteral: p.consume()})
	}
	return opt.Optional[NodeExpr]{}
}

func (p Parser) peek() opt.Optional[Token] {
	if p.currentIndex >= len(p.tokens) {
		return opt.NewOptional[Token]()
	}
	return opt.ToOptional(p.tokens[p.currentIndex])
}

func (p *Parser) consume() Token {
	r := p.tokens[p.currentIndex]
	p.currentIndex++
	return r
}
