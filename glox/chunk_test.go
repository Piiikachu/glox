package glox

import "testing"

func TestChunk(t *testing.T) {
	chunk := new(Chunk)
	constant := chunk.addConstant(Value{
		valueType: VAL_NUMBER,
		val:       12.3,
	})
	chunk.write(byte(OP_CONSTANT), 1)
	chunk.write(byte(constant), 1)
	chunk.write(byte(OP_RETURN), 1)
	disassemble(chunk, "test chunk")
}
