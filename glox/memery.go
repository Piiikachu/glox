package glox

func GROW_CAPACITY(capacity int) int {
	if capacity < 8 {
		return 8
	} else {
		return capacity * 2
	}
}

// func GROW_ARRAY(t reflect.Type,) []OpCode {

// }
