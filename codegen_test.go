package main

import (
	"bytes"
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
	loadAST := func(dir, filename string) *ASTNode {
		b, err := ioutil.ReadFile(pathToGoldenfile(dir, filename, ".ast"))
		if err != nil {
			t.Fatal(err)
		}
		program := new(ASTNode)
		err = json.Unmarshal(b, program)
		if err != nil {
			t.Fatal(err)
		}
		return program
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
				program := loadAST(dir, file.Name())
				dst := bytes.NewBuffer(make([]byte, 0, 4096))
				Compile(dst, program)
				if *viewAsm {
					fmt.Println(string(dst.Bytes()))
					return
				}
				goldenfilePath := pathToGoldenfile(dir, file.Name(), ".s")
				if *updateAsm {
					ioutil.WriteFile(goldenfilePath, dst.Bytes(), 0644)
					return
				}
				want, err := ioutil.ReadFile(goldenfilePath)
				if err != nil {
					t.Fatal(err)
				}
				if got := dst.Bytes(); !bytes.Equal(want, got) {
					t.Errorf("Expected:\n%s\nGot\n%s\n", want, got)
				}
			})
		}
	}
	runTests("tests/valid/")
}
