package main

import (
	"errors"
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
	identifier
	_var
	equals
	plus
	asterisk
	minus
	fslash
	openCurlyBracket
	closeCurlyBracket
	_if
)

func (t TokenType) GetBinPrec() opt.Optional[int] {
	switch t {
	case plus, minus:
		return opt.ToOptional(0)
	case asterisk, fslash:
		return opt.ToOptional(1)
	default:
		return opt.Optional[int]{}
	}
}

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

		} else if t.peek().MustGetValue() == '{' {
			t.consume()
			tokens = append(tokens, Token{tokenType: openCurlyBracket})

		} else if t.peek().MustGetValue() == '}' {
			t.consume()
			tokens = append(tokens, Token{tokenType: closeCurlyBracket})

		} else if t.peek().MustGetValue() == '=' {
			t.consume()
			tokens = append(tokens, Token{tokenType: equals})

		} else if t.peek().MustGetValue() == '+' {
			t.consume()
			tokens = append(tokens, Token{tokenType: plus})

		} else if t.peek().MustGetValue() == '*' {
			t.consume()
			tokens = append(tokens, Token{tokenType: asterisk})

		} else if t.peek().MustGetValue() == '-' {
			t.consume()
			tokens = append(tokens, Token{tokenType: minus})

		} else if t.peek().MustGetValue() == '/' {
			t.consume()
			if t.peek().HasValue() && t.peek().MustGetValue() == '/' {
				t.consume()
				for t.peek().HasValue() && t.peek().MustGetValue() != '\n' {
					t.consume()
				}
			} else if t.peek().HasValue() && t.peek().MustGetValue() == '*' {
				t.consume()
				for {
					if !t.peek().HasValue() {
						panic(errors.New("multiline comment didn't have an end. terminate it with `*/`"))
					}

					c := t.consume()
					if c == '*' && t.peek().HasValue() && t.peek().MustGetValue() == '/' {
						t.consume()
						break
					}
				}

			} else {
				tokens = append(tokens, Token{tokenType: fslash})
			}

		} else if unicode.IsLetter(t.peek().MustGetValue()) || t.peek().MustGetValue() == '$' || t.peek().MustGetValue() == '_' {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && (unicode.IsLetter(t.peek().MustGetValue()) || unicode.IsDigit(t.peek().MustGetValue()) || t.peek().MustGetValue() == '$' || t.peek().MustGetValue() == '_') {
				buf = append(buf, t.consume())
			}

			if string(buf) == "exit" {
				tokens = append(tokens, Token{tokenType: exit})
				buf = []rune{}
			} else if string(buf) == "var" {
				tokens = append(tokens, Token{tokenType: _var})
				buf = []rune{}
			} else if string(buf) == "if" {
				tokens = append(tokens, Token{tokenType: _if})
				buf = []rune{}
			} else {
				tokens = append(tokens, Token{tokenType: identifier, value: opt.ToOptional(string(buf))})
				buf = []rune{}
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
