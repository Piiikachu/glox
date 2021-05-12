package main

import "fmt"

func (c *Chunk) disassemble(name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < c.count; {
		offset = c.disassembleInstruction(offset)
	}
}

func (c *Chunk) disassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)
	instruction := c.code[offset]
	switch instruction {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
