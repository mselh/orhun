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
	tokens := make([]string, 0)
	start := 0
	skip := 0

	// string
	prog := []rune(string(text))
	for i, r := range prog {
		// search for cases
		if skip > 0 {
			skip--
			continue
		}

		if r == '\n' {
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		// numeric literal
		if unicode.IsDigit(r) {
			start := i
			fin := i + 1
			for unicode.IsDigit(prog[fin]) {
				fin++
			}
			tokens = append(tokens, string(prog[start:fin]))
			skip = fin - start - 1
			continue
		}

		if isIdent(r) {
			fmt.Println("start:", start)
			start = i
			fin := i + 1
			for isIdent2(prog[fin]) {
				fin++
			}
			tokens = append(tokens, string(prog[start:fin]))
			fmt.Println("end:", fin)
			skip = fin - start - 1
			continue
		}

		// last case is for punctuators
		p := string(prog[i : i+2])
		if punctLen := readPunct(p); punctLen > 0 {
			tokens = append(tokens, string(prog[i:i+punctLen]))
			skip = punctLen - 1
			continue
		}

		fmt.Println(r, "invalid token")

		fmt.Println(string(r), i)
	}

	for i, v := range tokens {
		fmt.Println(i, "'"+v+"'")
	}

	// parse
	type Node struct {
		lhs *Node // left hand side
		rhs *Node // right hand side

		// Block, used if kind == ND_BLOCK
		body *Node
	}

	type prog struct {
		body   *Node
		locals *Obj
	}

}
