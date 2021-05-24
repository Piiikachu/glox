package glox

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	scanner := new(Scanner)
	tokens := []Token{}
	src := []string{"and", "class", "else", "false", "for", "fun", "if", "nil", "or", "print", "return", "super", "this", "true", "var", "where"}
	for _, word := range src {
		scanner.init(word)
		tokens = append(tokens, scanner.scanToken())
	}
	fmt.Println(tokens)
}

func TestScanChars(t *testing.T) {
	scanner := new(Scanner)
	tokens := []Token{}
	src:=[]string{"(",")","{","}",",",".","-","+",";","/","*","!","!=","=","==","<","<=",">",">=","\"abc\""}
	for _, word := range src {
		scanner.init(word)
		tokens = append(tokens, scanner.scanToken())
	}
	fmt.Println(tokens)
}
