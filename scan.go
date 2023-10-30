package main

import (
	"fmt"
	"log"
	"unicode"
)

type token struct {
	kind string
	val  string
	pos  int
}

func tokenize(prog []rune) []token {
	tokens := make([]token, 0)
	start := 0
	skip := 0

	for i, r := range prog {
		// search for cases
		if skip > 0 {
			skip--
			continue
		}

		if string(r) == "/" && len(prog) > i && string(prog[i+1]) == "/" {
			fin := i + 1
			for string(prog[fin]) != "\n" {
				fin++
			}
			skip = fin
			continue

		}

		if r == '\n' {
			tokens = append(tokens, token{
				val:  string("\n"),
				pos:  i,
				kind: "EOL",
			})
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
			tokens = append(tokens, token{
				val:  string(prog[start:fin]),
				pos:  i,
				kind: "NUM",
			})
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
			tokens = append(tokens, token{
				val:  string(prog[start:fin]),
				pos:  i,
				kind: "IDENTIFIER",
			})
			fmt.Println("end:", fin)
			skip = fin - start - 1
			continue
		}

		// last case is for punctuators
		p := string(prog[i : i+2])
		if punctLen := readPunct(p); punctLen > 0 {
			tokens = append(tokens, token{
				val:  string(prog[i : i+punctLen]),
				pos:  i,
				kind: "PUNCT",
			})
			skip = punctLen - 1
			continue
		}

		log.Fatal(r, "invalid token")
	}

	return tokens

}
