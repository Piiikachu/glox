package glox

import "testing"

func TestChunk(t *testing.T) {
	chunk := new(Chunk)
	chunk.writeConstant(12.34, 1)
	chunk.write(byte(OP_RETURN), 1)
	disassemble(chunk, "test chunk")
}

func (c *Chunk) writeConstant(v interface{}, line int) {
	var constant int
	switch t := v.(type) {
	case float64:
		constant = c.addConstant(Value{
			valueType: VAL_NUMBER,
			val:       t,
		})
	case bool:
		constant = c.addConstant(Value{
			valueType: VAL_BOOL,
			val:       t,
		})
	}

	c.write(byte(OP_CONSTANT), line)
	c.write(byte(constant), line)
}
