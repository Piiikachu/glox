package main

func main() {
	chunk := new(Chunk)
	// c.initChunk()

	chunk.write(OP_RETURN)
	chunk.disassemble("test chunk")
	chunk.free()

}
