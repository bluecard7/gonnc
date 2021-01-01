package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	viewAsm   = flag.Bool("v-asm", false, "print generated assembly")
	updateAsm = flag.Bool("u-asm", false, "updates assembly file output")
)

func TestCompile(t *testing.T) {
	loadAST := func(dir, filename string, program *ASTNode) {
		b, err := ioutil.ReadFile(jsonFilepath(dir, filename))
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(b, program)
		if err != nil {
			t.Fatal(err)
		}
	}
	runTests := func(dir string) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			t.Run(dir+file.Name(), func(t *testing.T) {
				program := new(ASTNode)
				loadAST(dir, file.Name(), program)
				fmt.Println(string(program.JSON()))
				/*
					//tmpfile, err := ioutil.TempFile("", "*printed")
					r, w, err := os.Pipe()
					if err != nil {
						t.Fatal(err)
					}
					os.Stdout = w
					//defer os.Remove(tmpfile.Name())
					os.Pipe(tmpfile, os.Stdout,)
					Compile(program)
				*/
			})
		}
	}
	runTests("tests/valid/")
}
