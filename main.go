package main

import "fmt"

func main() {
	c := &Chunk{}
	c.initChunk()
	fmt.Printf("count:%d\tcapacity:%d\tcode:%v\n", c.count, c.capacity, c.code)
	for i := 0; i < 17; i++ {
		c.writeChunk(OP_TEST)
	}
	fmt.Printf("count:%d\tcapacity:%d\tcode:%v\n", c.count, c.capacity, c.code)

}
