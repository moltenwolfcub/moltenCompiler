package main

import (
	"errors"
	"fmt"
	"os"
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

func tokenise(str string) ([]Token, error) {
	tokens := []Token{}
	buf := ""

	for i := 0; i < len(str); i++ {
		r := rune(str[i])

		if unicode.IsSpace(r) {
			continue
		} else if unicode.IsLetter(r) {
			buf += string(str[i])
			i++
			for unicode.IsDigit(rune(str[i])) || unicode.IsLetter(rune(str[i])) {
				buf += string(str[i])
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
			buf += string(str[i])
			i++
			for unicode.IsDigit(rune(str[i])) {
				buf += string(str[i])
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

func assembleFromTokens(tokens []Token) string {
	output := "global _start\n_start:\n"
	for i, token := range tokens {
		if token.tokenType == exit {
			if i+1 < len(tokens) && tokens[i+1].tokenType == intLiteral {
				if i+2 < len(tokens) && tokens[i+2].tokenType == semiColon {
					output += "\tmov rax, 60\n"
					output += "\tmov rdi, " + tokens[i+1].value + "\n"
					output += "\tsyscall\n"
				}
			}
		}
	}
	return output
}

func main() {
	err := checkCLA()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	program, err := loadProgram(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tokens, err := tokenise(program)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	println(assembleFromTokens(tokens))
}

func checkCLA() error {
	if len(os.Args) != 2 {
		return errors.New("bad usage. correct usage is:\n\"molten <main.mltn>\"")
	}
	return nil
}

func loadProgram(fileName string) (string, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return "", errors.New("cannot read the file")
	}

	return string(file), nil
}
