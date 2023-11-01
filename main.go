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

	print(program)
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
