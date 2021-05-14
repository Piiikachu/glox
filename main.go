package main

func main() {

	vm.init()

	chunk := new(Chunk)
	// c.initChunk()
	constant := chunk.addConstant(2.86)
	chunk.write(byte(OP_CONSTANT), 123)
	chunk.write(byte(constant), 123)
	chunk.write(byte(OP_NEGATE), 123)
	chunk.write(byte(OP_RETURN), 123)
	disassemble(chunk, "test chunk")
	
	interpret(chunk)
	vm.free()
	chunk.free()

}
