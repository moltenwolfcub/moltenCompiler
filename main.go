package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Bad usage! Correct usage is:\n\"molten <main.mltn>\"")
		return
	}
	fileName := os.Args[1]

	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Cannot read the file")
		return
	}

	program := string(file)

	print(program)
}
