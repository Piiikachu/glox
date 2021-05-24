package main

import (
	"os"
	"glox/glox"
)

func main() {
	vm:=new(glox.VM)
	vm.Init()

	if len(os.Args) == 1 {
		vm.Repl()
	} else if len(os.Args) == 2 {
		vm.RunFile(os.Args[1])
	} else {
		os.Stderr.WriteString("Usage: glox [path]\n")
		os.Exit(64)
	}

	vm.Free()
}
