package glox

import "testing"

func TestVm(t *testing.T) {
	chunk := new(Chunk)
	chunk.writeConstant(1.2, 1)
	chunk.writeConstant(3.4, 1)
	chunk.write(byte(OP_ADD), 1)

	chunk.writeConstant(5.6,1)
	chunk.write(byte(OP_DIVIDE),1)

	chunk.write(byte(OP_NEGATE),1)

	chunk.write(byte(OP_RETURN), 1)
	disassemble(chunk, "test chunk")
}
