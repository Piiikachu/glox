package main

type OpCode byte

const (
	OP_CONSTANT OpCode = iota
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_ADD
	OP_SUBSTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN
)

type Chunk struct {
	count       int
	capacity    int
	code        []byte
	lines       []int
	currentCode int
	constants   ValueArray
}

func (c *Chunk) write(b byte, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
	c.capacity = cap(c.code)
	c.count = len(c.code)
}

func (c *Chunk) free() {
	c = new(Chunk)
}

func (c *Chunk) addConstant(value Value) int {
	c.constants.write(value)
	return c.constants.count - 1
}

func (c *Chunk) getCode() byte {
	return c.code[c.currentCode]
}
