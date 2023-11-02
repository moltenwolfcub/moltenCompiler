package main

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	exit TokenType = iota
	intLiteral
	semiColon
)

type Token struct {
	tokenType TokenType
	value     string
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

	for t.peek() >= 0 {
		if unicode.IsSpace(t.peek()) {
			t.consume()

		} else if t.peek() == ';' {
			t.consume()
			tokens = append(tokens, Token{tokenType: semiColon})

		} else if unicode.IsLetter(t.peek()) {
			buf = append(buf, t.consume())
			for t.peek() >= 0 && (unicode.IsLetter(t.peek()) || unicode.IsDigit(t.peek())) {
				buf = append(buf, t.consume())
			}

			if string(buf) == "exit" {
				tokens = append(tokens, Token{tokenType: exit})
				buf = []rune{}
			} else {
				return nil, fmt.Errorf("unknown keyword: %s", string(buf))
			}

		} else if unicode.IsDigit(t.peek()) {
			buf = append(buf, t.consume())
			for t.peek() >= 0 && unicode.IsDigit(t.peek()) {
				buf = append(buf, t.consume())
			}

			tokens = append(tokens, Token{tokenType: intLiteral, value: string(buf)})
			buf = []rune{}

		} else {
			return nil, fmt.Errorf("unknown token: %c", t.peek())
		}
	}

	return tokens, nil
}

func (t Tokeniser) peek() rune {
	if t.currentIndex >= len(t.program) {
		return -1
	}
	return rune(t.program[t.currentIndex])
}

func (t *Tokeniser) consume() rune {
	r := rune(t.program[t.currentIndex])
	t.currentIndex++
	return r
}
