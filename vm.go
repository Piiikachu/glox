package main

import "fmt"

const STACK_MAX = 256

var vm = new(VM)

type VM struct {
	chunk     Chunk
	ips       []byte
	stack     [STACK_MAX]Value
	stackTop  int
	currentIP int
}

type InterpretResult byte

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

func (vm *VM) init() {
	resetStack()
}

func resetStack() {
	vm.stackTop = 0
}

func (vm *VM) free() {

}

func (vm *VM) push(value Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}
func (vm *VM) pop() Value {
	vm.stackTop--
	value := vm.stack[vm.stackTop]
	return value
}

func interpret(c *Chunk) InterpretResult {
	vm.chunk = *c
	vm.ips = (*c).code
	return run()
}

func run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			fmt.Printf("          ")
			fmt.Println(vm.stack[:vm.stackTop])

			chunk := vm.getChunk()
			disassembleInstruction(chunk, vm.currentIP-chunk.currentCode)
		}

		instruction := OpCode(vm.READ_BYTE())
		switch instruction {
		case OP_CONSTANT:
			{
				constant := vm.READ_CONSTANT()
				vm.push(constant)
			}
		case OP_ADD:
			BINARY_OP('+')
		case OP_SUBSTRACT:
			BINARY_OP('-')
		case OP_MULTIPLY:
			BINARY_OP('*')
		case OP_DIVIDE:
			BINARY_OP('/')
		case OP_NEGATE:
			vm.push(-vm.pop())
		case OP_RETURN:
			{
				printValue(vm.pop())
				fmt.Println()
				return INTERPRET_OK
			}
		}
	}
}

func (vm *VM) READ_BYTE() byte {
	code := vm.getIP()
	vm.currentIP++
	return code
}

func (vm *VM) READ_CONSTANT() Value {
	return vm.getChunk().constants.values[vm.READ_BYTE()]
}

func (vm *VM) getChunk() *Chunk {
	return &vm.chunk
}

func (vm *VM) getIP() byte {
	return vm.ips[vm.currentIP]
}

func BINARY_OP(op rune) {
	b := vm.pop()
	a := vm.pop()
	switch op {
	case '+':
		vm.push(a + b)
	case '-':
		vm.push(a - b)
	case '*':
		vm.push(a * b)
	case '/':
		vm.push(a / b)
	}
}
