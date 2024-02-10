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

	file string
	line int
	col  int
}

type Tokeniser struct {
	program      string
	currentIndex int

	fileName  string
	lineCount int
	colCount  int
}

func NewTokeniser(program string, fileName string) Tokeniser {
	return Tokeniser{
		program:   program,
		fileName:  fileName,
		lineCount: 1,
		colCount:  1,
	}
}

func (t *Tokeniser) Tokenise() ([]Token, error) {
	tokens := []Token{}
	buf := []rune{}

	for t.peek().HasValue() {
		if t.peek().MustGetValue() == '\n' {
			t.consume()
			t.lineCount++
			t.colCount = 1

		} else if unicode.IsSpace(t.peek().MustGetValue()) {
			t.consume()
			t.colCount += 1

		} else if t.peek().MustGetValue() == ';' {
			t.consume()
			tokens = append(tokens, Token{tokenType: semiColon, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '(' {
			t.consume()
			tokens = append(tokens, Token{tokenType: openRoundBracket, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == ')' {
			t.consume()
			tokens = append(tokens, Token{tokenType: closeRoundBracket, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '{' {
			t.consume()
			tokens = append(tokens, Token{tokenType: openCurlyBracket, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '}' {
			t.consume()
			tokens = append(tokens, Token{tokenType: closeCurlyBracket, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '=' {
			t.consume()
			tokens = append(tokens, Token{tokenType: equals, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '+' {
			t.consume()
			tokens = append(tokens, Token{tokenType: plus, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '*' {
			t.consume()
			tokens = append(tokens, Token{tokenType: asterisk, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '-' {
			t.consume()
			tokens = append(tokens, Token{tokenType: minus, line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += 1

		} else if t.peek().MustGetValue() == '/' {
			t.consume()
			if t.peek().HasValue() && t.peek().MustGetValue() == '/' {
				t.consume()
				t.colCount += 1
				for t.peek().HasValue() && t.peek().MustGetValue() != '\n' {
					t.consume()
					t.colCount += 1
				}
			} else if t.peek().HasValue() && t.peek().MustGetValue() == '*' {
				t.consume()
				t.colCount += 1
				for {
					if !t.peek().HasValue() {
						return nil, t.error("multiline comment wasn't closed. terminate it with `*/`")
					}

					c := t.consume()
					t.colCount += 1
					if c == '*' && t.peek().HasValue() && t.peek().MustGetValue() == '/' {
						t.consume()
						t.colCount += 1
						break
					} else if c == '\n' {
						t.lineCount++
						t.colCount = 1
					}
				}

			} else {
				tokens = append(tokens, Token{tokenType: fslash, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += 1
			}

		} else if unicode.IsLetter(t.peek().MustGetValue()) || t.peek().MustGetValue() == '$' || t.peek().MustGetValue() == '_' {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && (unicode.IsLetter(t.peek().MustGetValue()) || unicode.IsDigit(t.peek().MustGetValue()) || t.peek().MustGetValue() == '$' || t.peek().MustGetValue() == '_') {
				buf = append(buf, t.consume())
			}

			if string(buf) == "exit" {
				tokens = append(tokens, Token{tokenType: exit, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else if string(buf) == "var" {
				tokens = append(tokens, Token{tokenType: _var, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else if string(buf) == "if" {
				tokens = append(tokens, Token{tokenType: _if, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else if string(buf) == "while" {
				tokens = append(tokens, Token{tokenType: while, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else if string(buf) == "else" {
				tokens = append(tokens, Token{tokenType: _else, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else if string(buf) == "break" {
				tokens = append(tokens, Token{tokenType: _break, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else if string(buf) == "continue" {
				tokens = append(tokens, Token{tokenType: _continue, line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			} else {
				tokens = append(tokens, Token{tokenType: identifier, value: opt.ToOptional(string(buf)), line: t.lineCount, col: t.colCount, file: t.fileName})
				t.colCount += len(buf)
				buf = []rune{}
			}

		} else if unicode.IsDigit(t.peek().MustGetValue()) {
			buf = append(buf, t.consume())
			for t.peek().HasValue() && unicode.IsDigit(t.peek().MustGetValue()) {
				buf = append(buf, t.consume())
			}

			tokens = append(tokens, Token{tokenType: intLiteral, value: opt.ToOptional(string(buf)), line: t.lineCount, col: t.colCount, file: t.fileName})
			t.colCount += len(buf)
			buf = []rune{}

		} else {
			return nil, t.error(fmt.Sprintf("invalid token: %c", t.peek().MustGetValue()))
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

func (t Tokeniser) error(message string) error {
	return fmt.Errorf("%s:%d:%d: %s", t.fileName, t.lineCount, t.colCount, message)
}
