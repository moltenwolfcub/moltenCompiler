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
			if !p.peek(1).HasValue() || p.peek(1).MustGetValue().tokenType != openRoundBracket {
				return opt.NewOptional[NodeExit](), errors.New("expected '(' after 'exit'")
			}

			p.consume()
			p.consume()

			if nodeExpr := p.ParseExpr(); nodeExpr.HasValue() {
				node.SetValue(NodeExit{expression: nodeExpr.MustGetValue()})
			} else if p.peek().HasValue() && p.peek().MustGetValue().tokenType == closeRoundBracket {
				node.SetValue(NodeExit{expression: NodeExpr{intLiteral: Token{tokenType: intLiteral, value: opt.ToOptional("0")}}})
			} else {
				return opt.NewOptional[NodeExit](), errors.New("invalid expression")
			}

			if p.peek().HasValue() && p.peek().MustGetValue().tokenType == closeRoundBracket {
				p.consume()
			} else {
				return opt.NewOptional[NodeExit](), errors.New("missing ')'")
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

func (p *Parser) ParseExpr() opt.Optional[NodeExpr] {
	if p.peek().HasValue() && p.peek().MustGetValue().tokenType == intLiteral {
		return opt.ToOptional(NodeExpr{intLiteral: p.consume()})
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
