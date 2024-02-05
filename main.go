package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("reading the program")
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
		fmt.Println(i, "'"+v.val+"'")
	}

	// parse
	//fn := new(Function)
	//fn.body = compoundStmt(tokens, fn)
	//fmt.Println(fn)

	// codegen

}
