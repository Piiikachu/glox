package main

import (
	"os"
)

func main() {

	vm.init()

	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		os.Stderr.WriteString("Usage: glox [path]\n")
		os.Exit(64)
	}

	vm.free()
}
