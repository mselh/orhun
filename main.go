package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var DEBUG = false

func myPrintln(a ...any) {
	if DEBUG {
		fmt.Println(a)
	}
}

func myPrint(a ...any) {
	if DEBUG {
		fmt.Print(a)
	}
}

func main() {
	myPrintln("reading the program")
	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Println(err)
	}

	// tokenize
	prog := []rune(string(text))
	r := reader{
		text: prog,
	}
	r.tokenize()
	for i, v := range r.tokenList {
		myPrintln(i, "'"+v.val+"'", "kind:", v.kind)
	}

	// parse
	p := parser{
		tokenList: r.tokenList,
		cur:       0,
		root: &node{
			kind: "root",
		},
	}
	p.parseAll()

	// walk
	walker := newProgram(p.root)
	walker.walk()

}
