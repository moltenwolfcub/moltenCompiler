package main

import (
	"fmt"
	"unicode"

	opt "github.com/moltenwolfcub/moltenCompiler/optional"
)

type TokenType int

const (
	exit TokenType = iota
	intLiteral
	semiColon
	openRoundBracket
	closeRoundBracket
)

type Token struct {
	tokenType TokenType
	value     opt.Optional[string]
}

type Tokeniser struct {
	program      string
	currentIndex int
}

func NewTokeniser(program string) Tokeniser {
	return Tokeniser{
		program: program,
	}
}

func (t Tokeniser) Tokenise() ([]Token, error) {
	tokens := []Token{}
	buf := []rune{}

	for t.peek().HasValue() {
		if unicode.IsSpace(t.peek().MustGetValue()) {
			t.consume()

		} else if t.peek().MustGetValue() == ';' {
			t.consume()
			tokens = append(tokens, Token{tokenType: semiColon})

		} else if t.peek().MustGetValue() == '(' {
			t.consume()
			tokens = append(tokens, Token{tokenType: openRoundBracket})

		} else if t.peek().MustGetValue() == ')' {
			t.consume()
			tokens = append(tokens, Token{tokenType: closeRoundBracket})

		} else if unicode.IsLetter(t.peek().MustGetValue()) {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && (unicode.IsLetter(t.peek().MustGetValue()) || unicode.IsDigit(t.peek().MustGetValue())) {
				buf = append(buf, t.consume())
			}

			if string(buf) == "exit" {
				tokens = append(tokens, Token{tokenType: exit})
				buf = []rune{}
			} else {
				return nil, fmt.Errorf("unknown keyword: %s", string(buf))
			}

		} else if unicode.IsDigit(t.peek().MustGetValue()) {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && unicode.IsDigit(t.peek().MustGetValue()) {
				buf = append(buf, t.consume())
			}

			tokens = append(tokens, Token{tokenType: intLiteral, value: opt.ToOptional(string(buf))})
			buf = []rune{}

		} else {
			return nil, fmt.Errorf("unknown token: %c", t.peek().MustGetValue())
		}
	}

	return tokens, nil
}

func (t Tokeniser) peek() opt.Optional[rune] {
	if t.currentIndex >= len(t.program) {
		return opt.NewOptional[rune]()
	}
	return opt.ToOptional(rune(t.program[t.currentIndex]))
}

func (t *Tokeniser) consume() rune {
	r := rune(t.program[t.currentIndex])
	t.currentIndex++
	return r
}
