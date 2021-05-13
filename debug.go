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
	switch OpCode(instruction) {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
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

func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.constants.values[constant])
	fmt.Println()
	return offset + 2
}

func printValue(value Value) {
	fmt.Printf("%g", value)
}
