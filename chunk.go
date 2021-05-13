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
	constants ValueArray
}

func (c *Chunk) write(b byte) {
	c.code = append(c.code, b)
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
