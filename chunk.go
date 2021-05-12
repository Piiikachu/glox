package main

type OpCode int8

const (
	OP_RETURN OpCode = iota
	OP_TEST
)

type Chunk struct {
	count    int
	capacity int
	code     []OpCode
}

func (c *Chunk) write(byteCode OpCode) {
	c.code = append(c.code, byteCode)
	c.capacity = cap(c.code)
	c.count = len(c.code)
}

func (c *Chunk) free() {
	c = new(Chunk)
}
