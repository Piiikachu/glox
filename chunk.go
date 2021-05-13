package main

type OpCode byte

const (
	OP_CONSTANT OpCode = iota
	OP_RETURN
)

type Chunk struct {
	count     int
	capacity  int
	code      []byte
	lines     []int
	constants ValueArray
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