package main

import (
	"fmt"
	"strconv"
)

func ASTToASM(node *ASTNode) {
	fmt.Println("\t.globl main")
	for _, child := range node.children {
		switch child.Get("Type") {
		case "Function":
			fmt.Printf("%s:\n", child.Get("Name")) // Prop instead of Get
			for _, stmt := range child.children {
				// remove Type, since its implied that a Function has statements?
				switch stmt.Get("StmtType") {
				case "Return":
					retVal, _ := strconv.ParseInt(stmt.Get("Value"), 10, 64)
					fmt.Printf("\tmov $%v, %%rax\n", retVal)
					fmt.Println("\tret")
				}
			}
		}
	}
}
