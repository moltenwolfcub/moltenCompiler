package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// if true molten program will run els it will just compile to asm
var ShouldRun = true

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
	rootNode := parser.ParseProg()

	root, err := rootNode.GetValue()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	generator := NewGenerator(root)
	asm := generator.GenProg()

	err = writeToFile(strings.Split(os.Args[1], ".")[0], asm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if ShouldRun {
		err = run(strings.Split(os.Args[1], ".")[0])
		if err != nil {
			fmt.Println(err.Error())
		}
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

func writeToFile(filename string, asm string) error {
	buildDir := "build"

	if _, err := os.Stat(buildDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(buildDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating logger directory: %v", err.Error())
		}
	}

	file, err := os.Create(buildDir + "/" + filename + ".asm")
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

func run(filename string) error {
	cmd := exec.Command("./run.sh", filename+".asm")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
