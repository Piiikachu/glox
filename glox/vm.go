package glox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const STACK_MAX = 256

type VM struct {
	chunk     *Chunk
	ips       []byte
	stack     [STACK_MAX]Value
	stackTop  int
	currentIP int
	objects   Obj
}

type InterpretResult byte

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

func (vm *VM) Repl() {
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			os.Stderr.WriteString("Read error")
			return
		}
		vm.interpret(line)
	}
}

func (vm *VM) RunFile(path string) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	result := vm.interpret(string(buffer))

	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}
}

func (vm *VM) Init() {
	vm.resetStack()
	vm.objects = nil
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}

func (vm *VM) Free() {
	vm.resetStack()
	vm.freeObjects()
}

func (vm *VM) freeObjects() {
	obj := vm.objects
	for obj != nil && obj.next() != nil {
		next := obj.next()
		if next == nil {
			return
		}
		(*next).free()
		obj = *next
	}
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

func (vm *VM) interpret(source string) InterpretResult {
	chunk := new(Chunk)
	if !vm.compile(source, chunk) {
		return INTERPRET_COMPILE_ERROR
	}
	vm.chunk = chunk
	vm.ips = vm.chunk.code
	result := vm.run()
	return result
}

func (vm *VM) run() InterpretResult {
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
			vm.push(NIL_VAL())
		case OP_TRUE:
			vm.push(BOOL_VAL(true))
		case OP_FALSE:
			vm.push(BOOL_VAL(false))
		case OP_EQUAL:
			{
				b := vm.pop()
				a := vm.pop()
				vm.push(BOOL_VAL(a.equals(b)))
			}
		case OP_GREATER:
			vm.BINARY_OP('>')
		case OP_LESS:
			vm.BINARY_OP('<')
		case OP_ADD:
			v1 := vm.peek(0)
			v2 := vm.peek(1)
			if v1.isString() && v2.isString() {
				vm.concatenate()
			} else if v1.isType(VAL_NUMBER) && v2.isType(VAL_NUMBER) {
				b := vm.pop().asNumber()
				a := vm.pop().asNumber()
				vm.push(NUMBER_VAL(a + b))
			} else {
				vm.runtimeError("Operands must be two numbers or two strings.")
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_SUBSTRACT:
			vm.BINARY_OP('-')
		case OP_MULTIPLY:
			vm.BINARY_OP('*')
		case OP_DIVIDE:
			vm.BINARY_OP('/')
		case OP_NOT:
			vm.push(BOOL_VAL(isFalsey(vm.pop())))
		case OP_NEGATE:
			{
				if !vm.peek(0).isType(VAL_NUMBER) {
					vm.runtimeError("Operand must be a number.")
					return INTERPRET_RUNTIME_ERROR
				}
				vm.push(NUMBER_VAL(-vm.pop().asNumber()))
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

func (vm *VM) concatenate() {
	b := vm.pop().asString()
	a := vm.pop().asString()
	length := b.length + a.length
	str := a.str + b.str
	result := ObjString{
		length: length,
		str:    str,
	}
	vm.push(OBJ_VAL(&result))
}

func (vm *VM) runtimeError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a)

	instruction := vm.currentIP - vm.chunk.currentCode - 1
	line := vm.chunk.lines[instruction]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	vm.resetStack()
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

func (vm *VM) BINARY_OP(op rune) InterpretResult {
	if !vm.peek(0).isType(VAL_NUMBER) || !vm.peek(1).isType(VAL_NUMBER) {
		vm.runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}
	b := vm.pop().asNumber()
	a := vm.pop().asNumber()
	switch op {
	case '+':
		vm.push(NUMBER_VAL(a + b))
	case '-':
		vm.push(NUMBER_VAL(a - b))
	case '*':
		vm.push(NUMBER_VAL(a * b))
	case '/':
		vm.push(NUMBER_VAL(a / b))
	case '<':
		vm.push(BOOL_VAL(a < b))
	case '>':
		vm.push(BOOL_VAL(a > b))
	}
	return INTERPRET_OK
}
