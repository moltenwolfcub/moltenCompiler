package main

type Generator struct {
	root NodeExit
}

func NewGenerator(root NodeExit) Generator {
	return Generator{
		root: root,
	}
}

func (g Generator) Generate() string {
	output := "global _start\n_start:\n"

	output += "\tmov rax, 60\n"
	output += "\tmov rdi, " + g.root.expression.intLiteral.value.MustGetValue() + "\n"
	output += "\tsyscall\n"

	return output
}
