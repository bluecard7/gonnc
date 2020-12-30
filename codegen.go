package main

import (
	"fmt"
)

func Compile(program *ASTNode) {
	fmt.Println("\t.globl main")
	for _, child := range program.Children {
		astToAsm(child)
	}
}

func astToAsm(node *ASTNode) {
	switch node.Kind {
	case FUNCTION:
		// til functions are actually implemented
		fmt.Printf("%s:\n", "main")
		for _, child := range node.Children {
			astToAsm(child)
		}
	case RETURN:
		astToAsm(node.Children[0])
		fmt.Println("\tret")
	case NUM:
		fmt.Printf("\tmov $%d, %%rax\n", node.Value)
	case MUL:
		setupBinOp(node)
		fmt.Println("\timul %rdi, %rax")
	case DIV:
		setupBinOp(node)
		fmt.Println("\tcqo")
		fmt.Println("\tidiv %rdi")
	case ADD:
		setupBinOp(node)
		fmt.Println("\tadd %rdi, %rax")
	case SUB:
		setupBinOp(node)
		fmt.Println("\tsub %rdi, %rax")
	case NEG:
		astToAsm(node.Children[0])
		fmt.Println("\tneg %rax")
	}
}

func setupBinOp(node *ASTNode) {
	astToAsm(node.Children[1])
	fmt.Println("\tpush %rax")
	astToAsm(node.Children[0])
	fmt.Println("\tpop %rdi")
}
