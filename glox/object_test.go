package glox

import "testing"

func TestObject(t *testing.T) {
	source := "\"test\"+\" adidas\""
	// source := "\"test\"==\" adidas\""
	result := interpret(source)
	if result != INTERPRET_OK {
		t.Errorf("Interpret failed: %s", source)
	}

}
