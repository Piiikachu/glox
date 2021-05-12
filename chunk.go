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

func (c *Chunk) initChunk() {
	c.count = 0
	c.capacity = 0
	c.code = nil
}

func (c *Chunk) writeChunk(byteCode OpCode) {
	//capacity not enough
	if c.capacity < c.count+1 {
		oldCapacity := c.capacity
		c.capacity = GROW_CAPACITY(oldCapacity)
		c.code = make([]OpCode, c.capacity)
	}

	c.code[c.count] = byteCode
	c.count++
}
