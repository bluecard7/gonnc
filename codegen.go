package main

import (
	"fmt"
	"io"
)

func Compile(w io.Writer, program *ASTNode) {
	fmt.Fprintln(w, "\t.globl main")
	for _, child := range program.Children {
		astToAsm(w, child)
	}
}

func astToAsm(w io.Writer, node *ASTNode) {
	switch node.Kind {
	case FUNCTION:
		// til functions are actually implemented
		fmt.Fprintf(w, "%s:\n", "main")
		for _, child := range node.Children {
			astToAsm(w, child)
		}
	case RETURN:
		astToAsm(w, node.Children[0])
		fmt.Fprintln(w, "\tret")
	case NUM:
		fmt.Fprintf(w, "\tmov $%d, %%rax\n", node.Value)
	case MUL:
		setupBinOp(w, node)
		fmt.Fprintf(w, "\timul %%rdi, %%rax\n")
	case DIV:
		setupBinOp(w, node)
		fmt.Fprintf(w, "\tcqo\n")
		fmt.Fprintf(w, "\tidiv %%rdi\n")
	case ADD:
		setupBinOp(w, node)
		fmt.Fprintf(w, "\tadd %%rdi, %%rax\n")
	case SUB:
		setupBinOp(w, node)
		fmt.Fprintf(w, "\tsub %%rdi, %%rax\n")
	case NEG:
		astToAsm(w, node.Children[0])
		fmt.Fprintf(w, "\tneg %%rax\n")
	}
}

func setupBinOp(w io.Writer, node *ASTNode) {
	astToAsm(w, node.Children[1])
	fmt.Fprintf(w, "\tpush %%rax\n")
	astToAsm(w, node.Children[0])
	fmt.Fprintf(w, "\tpop %%rdi\n")
}
