package glox

import "testing"

func TestObject(t *testing.T) {
	vm:=new(VM)
	vm.Init()
	source := "\"test\"+\" adidas\""
	// source := "\"test\"==\" adidas\""
	result := vm.interpret(source)
	if result != INTERPRET_OK {
		t.Errorf("Interpret failed: %s", source)
	}
	vm.Free()
}
