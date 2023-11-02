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
	program string
}

func NewTokeniser(program string) Tokeniser {
	return Tokeniser{
		program: program,
	}
}

func (t Tokeniser) Tokenise() ([]Token, error) {
	tokens := []Token{}
	buf := ""

	for i := 0; i < len(t.program); i++ {
		r := rune(t.program[i])

		if unicode.IsSpace(r) {
			continue
		} else if unicode.IsLetter(r) {
			buf += string(t.program[i])
			i++
			for unicode.IsDigit(rune(t.program[i])) || unicode.IsLetter(rune(t.program[i])) {
				buf += string(t.program[i])
				i++
			}
			i--

			if buf == "exit" {
				tokens = append(tokens, Token{tokenType: exit})

				buf = ""
			} else {
				return nil, fmt.Errorf("unknown keyword: %s", buf)
			}
		} else if unicode.IsDigit(r) {
			buf += string(t.program[i])
			i++
			for unicode.IsDigit(rune(t.program[i])) {
				buf += string(t.program[i])
				i++
			}
			i--
			tokens = append(tokens, Token{tokenType: intLiteral, value: buf})

			buf = ""
		} else if r == ';' {
			tokens = append(tokens, Token{tokenType: semiColon})
		} else {
			return nil, fmt.Errorf("unknown token: %c", r)
		}
	}

	return tokens, nil
}
