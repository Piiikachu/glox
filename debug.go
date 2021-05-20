package main

import "fmt"

func disassemble(c *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < c.count; {
		offset = disassembleInstruction(c, offset)
	}
}

func disassembleInstruction(c *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && c.lines[offset] == c.lines[offset-1] {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", c.lines[offset])
	}

	instruction := c.code[offset]
	switch OpCode(instruction) {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OP_NIL:
		return simpleInstruction("OP_NIL", offset)
	case OP_TRUE:
		return simpleInstruction("OP_TRUE", offset)
	case OP_FALSE:
		return simpleInstruction("OP_FALSE", offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBSTRACT:
		return simpleInstruction("OP_SUBSTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
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
