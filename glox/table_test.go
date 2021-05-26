package glox

import (
	"fmt"
	"testing"
)

func TestTable(t *testing.T) {
	vm = new(VM)
	vm.Init()
	fmt.Println(vm.strings)
	newObjString("test")
	fmt.Println(vm.strings)
	newObjString("adidas")
	fmt.Println(vm.strings)
	newObjString("test")
	fmt.Println(vm.strings)
	newObjString("test")
	fmt.Println(vm.strings)
}

func TestEqual(t *testing.T) {
	vm := new(VM)
	vm.Init()
	source := "\" adidas\"==\" adidas\""
	result := vm.interpret(source)
	if result != INTERPRET_OK {
		t.Errorf("Interpret failed: %s", source)
	}
	vm.Free()
}
