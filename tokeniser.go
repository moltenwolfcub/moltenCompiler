package main

import (
	"fmt"
	"unicode"

	opt "github.com/moltenwolfcub/moltenCompiler/optional"
)

type TokenType int

const (
	intLiteral TokenType = iota
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
	while
	_else
	_break
	_continue
	_func
	comma
	_return
	syscall
	ampersand
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

	lineInfo LineInfo
}

type Tokeniser struct {
	program      string
	currentIndex int

	currentLineInfo LineInfo
}

func NewTokeniser(program string, fileName string) Tokeniser {
	return Tokeniser{
		program:         program,
		currentLineInfo: NewLineInfo(fileName),
	}
}

func (t *Tokeniser) Tokenise() ([]Token, error) {
	tokens := []Token{}
	buf := []rune{}

	for t.peek().HasValue() {
		if t.peek().MustGetValue() == '\n' {
			t.consume()
			t.currentLineInfo.NextLine()

		} else if unicode.IsSpace(t.peek().MustGetValue()) {
			t.consume()
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == ';' {
			t.consume()
			tokens = append(tokens, Token{tokenType: semiColon, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '(' {
			t.consume()
			tokens = append(tokens, Token{tokenType: openRoundBracket, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == ')' {
			t.consume()
			tokens = append(tokens, Token{tokenType: closeRoundBracket, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '{' {
			t.consume()
			tokens = append(tokens, Token{tokenType: openCurlyBracket, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '}' {
			t.consume()
			tokens = append(tokens, Token{tokenType: closeCurlyBracket, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '=' {
			t.consume()
			tokens = append(tokens, Token{tokenType: equals, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '+' {
			t.consume()
			tokens = append(tokens, Token{tokenType: plus, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '*' {
			t.consume()
			tokens = append(tokens, Token{tokenType: asterisk, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '-' {
			t.consume()
			tokens = append(tokens, Token{tokenType: minus, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == ',' {
			t.consume()
			tokens = append(tokens, Token{tokenType: comma, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '&' {
			t.consume()
			tokens = append(tokens, Token{tokenType: ampersand, lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncColumn()

		} else if t.peek().MustGetValue() == '/' {
			t.consume()
			if t.peek().HasValue() && t.peek().MustGetValue() == '/' {
				t.consume()
				t.currentLineInfo.IncColumn()
				for t.peek().HasValue() && t.peek().MustGetValue() != '\n' {
					t.consume()
					t.currentLineInfo.IncColumn()
				}
			} else if t.peek().HasValue() && t.peek().MustGetValue() == '*' {
				t.consume()
				t.currentLineInfo.IncColumn()
				for {
					if !t.peek().HasValue() {
						return nil, t.currentLineInfo.PositionedError("multiline comment wasn't closed. terminate it with `*/`")
					}

					c := t.consume()
					t.currentLineInfo.IncColumn()
					if c == '*' && t.peek().HasValue() && t.peek().MustGetValue() == '/' {
						t.consume()
						t.currentLineInfo.IncColumn()
						break
					} else if c == '\n' {
						t.currentLineInfo.NextLine()
					}
				}

			} else {
				tokens = append(tokens, Token{tokenType: fslash, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncColumn()
			}

		} else if unicode.IsLetter(t.peek().MustGetValue()) || t.peek().MustGetValue() == '$' || t.peek().MustGetValue() == '_' {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && (unicode.IsLetter(t.peek().MustGetValue()) || unicode.IsDigit(t.peek().MustGetValue()) || t.peek().MustGetValue() == '$' || t.peek().MustGetValue() == '_') {
				buf = append(buf, t.consume())
			}

			if string(buf) == "var" {
				tokens = append(tokens, Token{tokenType: _var, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "if" {
				tokens = append(tokens, Token{tokenType: _if, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "while" {
				tokens = append(tokens, Token{tokenType: while, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "else" {
				tokens = append(tokens, Token{tokenType: _else, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "break" {
				tokens = append(tokens, Token{tokenType: _break, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "continue" {
				tokens = append(tokens, Token{tokenType: _continue, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "func" {
				tokens = append(tokens, Token{tokenType: _func, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "return" {
				tokens = append(tokens, Token{tokenType: _return, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else if string(buf) == "syscall" {
				tokens = append(tokens, Token{tokenType: syscall, lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			} else {
				tokens = append(tokens, Token{tokenType: identifier, value: opt.ToOptional(string(buf)), lineInfo: t.currentLineInfo})
				t.currentLineInfo.IncWord(buf)
				buf = []rune{}
			}

		} else if unicode.IsDigit(t.peek().MustGetValue()) {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && unicode.IsDigit(t.peek().MustGetValue()) {
				buf = append(buf, t.consume())
			}

			tokens = append(tokens, Token{tokenType: intLiteral, value: opt.ToOptional(string(buf)), lineInfo: t.currentLineInfo})
			t.currentLineInfo.IncWord(buf)
			buf = []rune{}

		} else {
			return nil, t.currentLineInfo.PositionedError(fmt.Sprintf("invalid token: %c", t.peek().MustGetValue()))
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
