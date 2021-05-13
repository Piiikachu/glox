package main

func main() {
	chunk := new(Chunk)
	// c.initChunk()
	constant:=chunk.addConstant(2.86)
	chunk.write(byte(OP_CONSTANT))
	chunk.write(byte(constant))
	chunk.write(byte(OP_RETURN))
	chunk.disassemble("test chunk")
	chunk.free()

}
