package main

func main() {
	chunk := new(Chunk)
	// c.initChunk()
	constant := chunk.addConstant(2.86)
	chunk.write(byte(OP_CONSTANT), 123)
	chunk.write(byte(constant), 123)
	chunk.write(byte(OP_RETURN), 123)
	chunk.disassemble("test chunk")
	chunk.free()

}
