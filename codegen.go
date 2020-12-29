package main

import (
	"fmt"
	"strconv"
)

func ASTToASM(node *ASTNode) {
	fmt.Println("\t.globl main")
	for _, child := range node.Children {
		switch child.Kind {
		case FUNCTION:
			// til functions are actually implemented
			fmt.Printf("%s:\n", "main")
			ASTToASM(child)
		case RETURN:
			retVal, _ := strconv.ParseInt("0", 10, 64)
			fmt.Printf("\tmov $%v, %%rax\n", retVal)
			fmt.Println("\tret")
		case MUL:
			fmt.Println("\timul %%rdi, %%rax")
		case DIV:
			fmt.Println("\tcqo")
			fmt.Println("\tidiv %%rdi")
		case ADD:
			fmt.Println("\tadd %%rdi, %%rax")
		case SUB:
			fmt.Println("\tsub %%rdi, %%rax")
		}
	}
}
