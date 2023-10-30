package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

func isIdent(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// Returns true if c is valid as a non-first character of an identifier.
func isIdent2(r rune) bool {
	return isIdent(r) || ('0' <= r && r <= '9')
}

// All punctuations in C:
// ! " # $ % & ' ( ) * + , - . / : ;
// < = > ? @ [ \ ] ^ _ ` { | } ~
func readPunct(c string) int {

	if strings.HasPrefix(c, "==") ||
		strings.HasPrefix(c, "!=") ||
		strings.HasPrefix(c, "<=") ||
		strings.HasPrefix(c, ">=") {
		return 2
	}

	if unicode.IsPunct(rune(c[0])) || strings.HasPrefix(c, "+") ||
		strings.HasPrefix(c, "-") ||
		strings.HasPrefix(c, "<") ||
		strings.HasPrefix(c, ">") ||
		strings.HasPrefix(c, "=") ||
		strings.HasPrefix(c, ";") {
		return 1
	}
	return 0
}

func main() {
	fmt.Println("reading the program")
	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Println(err)
	}

	// tokenize
	prog := []rune(string(text))
	tokens := tokenize(prog)
	for i, v := range tokens {
		fmt.Println(i, "'"+v.val+"'")
	}

	// parse
	fn := new(Function)
	fn.body = compoundStmt(tokens, fn)
	fmt.Println(fn)

	// codegen

}
