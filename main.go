package main

import (
	"errors"
	"fmt"
	"os"
)

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

	tokeniser := NewTokeniser(program)

	tokens, err := tokeniser.Tokenise()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	asm := assembleFromTokens(tokens)
	err = writeToFile(asm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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

func writeToFile(asm string) error {
	buildDir := "build"

	if _, err := os.Stat(buildDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(buildDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating logger directory: %v", err.Error())
		}
	}

	file, err := os.Create(buildDir + "/test.asm")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(asm)
	if err != nil {
		return err
	}
	return nil
}
