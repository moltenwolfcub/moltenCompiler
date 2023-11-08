package main

import (
	"errors"
	"fmt"
	"os"
)

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

	parser := NewParser(tokens)
	rootNode, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	root, err := rootNode.GetValue()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	generator := NewGenerator(root)
	asm := generator.Generate()

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
