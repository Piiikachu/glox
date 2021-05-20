package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const STACK_MAX = 256

var vm = new(VM)

type VM struct {
	chunk     *Chunk
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

func repl() {
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			os.Stderr.WriteString("Read error")
			return
		}
		interpret(line)
	}
}

func runFile(path string) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	result := interpret(string(buffer))

	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}
}

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

func interpret(source string) InterpretResult {
	chunk := new(Chunk)
	if !compile(source, chunk) {
		return INTERPRET_COMPILE_ERROR
	}
	vm.chunk = chunk
	vm.ips = vm.chunk.code
	result := run()
	return result
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
		case OP_NIL:
			vm.push(*VAL(nil))
		case OP_TRUE:
			vm.push(*VAL(true))
		case OP_FALSE:
			vm.push(*VAL(false))
		case OP_ADD:
			BINARY_OP('+')
		case OP_SUBSTRACT:
			BINARY_OP('-')
		case OP_MULTIPLY:
			BINARY_OP('*')
		case OP_DIVIDE:
			BINARY_OP('/')
		case OP_NOT:
			vm.push(*VAL(isFalsey(vm.pop())))
		case OP_NEGATE:
			{
				if !vm.peek(0).isType(VAL_NUMBER) {
					runtimeError("Operand must be a number.")
					return INTERPRET_RUNTIME_ERROR
				}
				vm.push(*VAL(-vm.pop().asNumber()))
			}
		case OP_RETURN:
			{
				printValue(vm.pop())
				fmt.Println()
				return INTERPRET_OK
			}
		}
	}
}

func (vm *VM) peek(offset int) Value {
	return vm.stack[vm.stackTop-1-offset]
}

func isFalsey(value Value) bool {
	//nil is falsey
	if value.isType(VAL_NIL) {
		return true
	} else if value.isType(VAL_BOOL) {
		//false is falsey
		if !value.asBool() {
			return true
		}
	}
	return false
}

func runtimeError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a)

	instruction := vm.currentIP - vm.chunk.currentCode - 1
	line := vm.chunk.lines[instruction]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	resetStack()
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
	return vm.chunk
}

func (vm *VM) getIP() byte {
	return vm.ips[vm.currentIP]
}

func BINARY_OP(op rune) InterpretResult {
	if !vm.peek(0).isType(VAL_NUMBER) || !vm.peek(1).isType(VAL_NUMBER) {
		runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}
	b := vm.pop().asNumber()
	a := vm.pop().asNumber()
	switch op {
	case '+':
		vm.push(*VAL(a + b))
	case '-':
		vm.push(*VAL(a - b))
	case '*':
		vm.push(*VAL(a * b))
	case '/':
		vm.push(*VAL(a / b))
	}
	return INTERPRET_OK
}
