package main

import "fmt"

func compile(source string) {
	scanner := new(Scanner)
	scanner.init(source)
	line := -1
	for {
		token := scanner.scanToken()
		if token.line != line {
			fmt.Printf("[%4d] ", token.line)
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%-20s '%.*s' \n", token.tokenType.String(), len(token.lexeme), token.lexeme)

		if token.tokenType == TOKEN_EOF {
			break
		}
	}

}
