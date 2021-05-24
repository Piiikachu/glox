package main

import (
	"fmt"
	"math"
	"os"
)

type Parser struct {
	current   Token
	previous  Token
	hadError  bool
	panicMode bool
}

type Precedence byte

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type ParseFn func()

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

var parser = &Parser{
	hadError:  false,
	panicMode: false,
}

var compilingChunk *Chunk

var rules map[TokenType]ParseRule

func compile(source string, chunk *Chunk) bool {
	scanner.init(source)
	compilingChunk = chunk
	parser.advance()
	parser.expression()
	parser.consume(TOKEN_EOF, "Expect end of expression.")
	parser.endCompiler()
	return !parser.hadError
}

func (p *Parser) advance() {
	p.previous = p.current
	for {
		p.current = scanner.scanToken()
		if p.current.tokenType != TOKEN_ERROR {
			break
		}
		errorAtCurrent(p.current.lexeme)
	}
}

func (p *Parser) expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}

func parsePrecedence(p Precedence) {
	parser.advance()
	prefix := rules[parser.previous.tokenType].prefix
	if prefix == nil {
		errorAtPrevious("Expect expression.")
		return
	}
	prefix()

	for p <= rules[parser.current.tokenType].precedence {
		parser.advance()
		infix := rules[parser.previous.tokenType].infix
		infix()
	}
}

func number() {
	value :=NUMBER_VAL(parser.previous.literal.(float64))
	emitConstant(value)
}

func grouping() {
	parser.expression()
	parser.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func unary() {
	operator := parser.previous.tokenType

	parsePrecedence(PREC_UNARY)

	switch operator {
	case TOKEN_BANG:
		emitByte(byte(OP_NOT))
	case TOKEN_MINUS:
		emitByte(byte(OP_NEGATE))
	default:
		return
	}
}

func binary() {
	operator := parser.previous.tokenType

	parsePrecedence(rules[operator].precedence + 1)

	switch operator {
	case TOKEN_BANG_EQUAL:
		emitBytes(byte(OP_EQUAL), byte(OP_NOT))
	case TOKEN_EQUAL_EQUAL:
		emitByte(byte(OP_EQUAL))
	case TOKEN_GREATER_EQUAL:
		emitBytes(byte(OP_LESS), byte(OP_NOT))
	case TOKEN_GREATER:
		emitByte(byte(OP_GREATER))
	case TOKEN_LESS_EQUAL:
		emitBytes(byte(OP_GREATER), byte(OP_NOT))
	case TOKEN_LESS:
		emitByte(byte(OP_LESS))
	case TOKEN_PLUS:
		emitByte(byte(OP_ADD))
	case TOKEN_MINUS:
		emitByte(byte(OP_SUBSTRACT))
	case TOKEN_STAR:
		emitByte(byte(OP_MULTIPLY))
	case TOKEN_SLASH:
		emitByte(byte(OP_DIVIDE))
	default:
		return
	}

}

func literal() {
	switch parser.previous.tokenType {
	case TOKEN_FALSE:
		emitByte(byte(OP_FALSE))
	case TOKEN_NIL:
		emitByte(byte(OP_NIL))
	case TOKEN_TRUE:
		emitByte(byte(OP_TRUE))
	default:
		return
	}
}

func gstring() {
	str := parser.previous.literal.(string)
	emitConstant(OBJ_VAL(newObjString(str)))
}

func (p *Parser) consume(t TokenType, msg string) {
	if parser.current.tokenType == t {
		p.advance()
		return
	}
	errorAtCurrent(msg)
}

func (p *Parser) endCompiler() {
	emitReturn()
	if DEBUG_PRINT_CODE {
		if !parser.hadError {
			disassemble(currentChunk(), "code")
		}
	}
}

func emitConstant(value Value) {
	emitBytes(byte(OP_CONSTANT), makeConstant(value))
}

func emitReturn() {
	emitByte(byte(OP_RETURN))
}

func emitByte(b byte) {
	currentChunk().write(b, parser.previous.line)
}

func emitBytes(b1 byte, b2 byte) {
	emitByte(b1)
	emitByte(b2)
}

func makeConstant(value Value) byte {
	constant := compilingChunk.addConstant(value)
	if constant > math.MaxUint8 {
		errorAtPrevious("Too many constants in one chunk.")
		return 0
	}
	return byte(constant)
}

func currentChunk() *Chunk {
	return compilingChunk
}

func errorAtCurrent(msg string) {
	errorAt(&parser.current, msg)
}

func errorAtPrevious(msg string) {
	errorAt(&parser.previous, msg)
}

func errorAt(token *Token, msg string) {
	if parser.hadError {
		return
	}
	parser.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)
	if token.tokenType == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.tokenType == TOKEN_ERROR {

	} else {
		fmt.Fprintf(os.Stderr, " at '%.*s'", len(token.lexeme), token.lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", msg)

	parser.hadError = true
}

func init() {
	rules = map[TokenType]ParseRule{
		TOKEN_LEFT_PAREN:    {grouping, nil, PREC_NONE},
		TOKEN_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		TOKEN_LEFT_BRACE:    {nil, nil, PREC_NONE},
		TOKEN_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		TOKEN_COMMA:         {nil, nil, PREC_NONE},
		TOKEN_DOT:           {nil, nil, PREC_NONE},
		TOKEN_MINUS:         {unary, binary, PREC_TERM},
		TOKEN_PLUS:          {nil, binary, PREC_TERM},
		TOKEN_SEMICOLON:     {nil, nil, PREC_NONE},
		TOKEN_SLASH:         {nil, binary, PREC_FACTOR},
		TOKEN_STAR:          {nil, binary, PREC_FACTOR},
		TOKEN_BANG:          {unary, nil, PREC_NONE},
		TOKEN_BANG_EQUAL:    {nil, binary, PREC_EQUALITY},
		TOKEN_EQUAL:         {nil, nil, PREC_NONE},
		TOKEN_EQUAL_EQUAL:   {nil, binary, PREC_EQUALITY},
		TOKEN_GREATER:       {nil, binary, PREC_COMPARISON},
		TOKEN_GREATER_EQUAL: {nil, binary, PREC_COMPARISON},
		TOKEN_LESS:          {nil, binary, PREC_COMPARISON},
		TOKEN_LESS_EQUAL:    {nil, binary, PREC_COMPARISON},
		TOKEN_IDENTIFIER:    {nil, nil, PREC_NONE},
		TOKEN_STRING:        {gstring, nil, PREC_NONE},
		TOKEN_NUMBER:        {number, nil, PREC_NONE},
		TOKEN_AND:           {nil, nil, PREC_NONE},
		TOKEN_CLASS:         {nil, nil, PREC_NONE},
		TOKEN_ELSE:          {nil, nil, PREC_NONE},
		TOKEN_FALSE:         {literal, nil, PREC_NONE},
		TOKEN_FOR:           {nil, nil, PREC_NONE},
		TOKEN_FUN:           {nil, nil, PREC_NONE},
		TOKEN_IF:            {nil, nil, PREC_NONE},
		TOKEN_NIL:           {literal, nil, PREC_NONE},
		TOKEN_OR:            {nil, nil, PREC_NONE},
		TOKEN_PRINT:         {nil, nil, PREC_NONE},
		TOKEN_RETURN:        {nil, nil, PREC_NONE},
		TOKEN_SUPER:         {nil, nil, PREC_NONE},
		TOKEN_THIS:          {nil, nil, PREC_NONE},
		TOKEN_TRUE:          {literal, nil, PREC_NONE},
		TOKEN_VAR:           {nil, nil, PREC_NONE},
		TOKEN_WHILE:         {nil, nil, PREC_NONE},
		TOKEN_ERROR:         {nil, nil, PREC_NONE},
		TOKEN_EOF:           {nil, nil, PREC_NONE},
	}

}
