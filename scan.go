package main

import (
	"fmt"
	"log"
	"regexp"
	"unicode"
)

type token struct {
	kind string
	val  string
	pos  int
	line int
}

type reader struct {
	text      []rune
	cur       int
	tokenList []token
	lineNo    int
}

func (r *reader) tokenize() {
	// dont token whitespaces
	for r.cur < len(r.text) {
		c := r.now()

		if r.now() == '/' && r.cur < len(r.text) && r.peekN(1) == '*' {
			fmt.Println("comment token!!!")

			fmt.Println("r back", string(r.backstep()), "bool:", r.backstep() != '*')
			fmt.Println("r now", string(r.now()), "bool:", r.now() != '/')
			// skip reading until */
			for {
				if r.cur == len(r.text) {
					log.Fatalln("comment is not closed")
				}
				if r.now() == '\n' {
					r.lineNo++
				}
				if r.now() == '/' && r.backstep() == '*' {
					r.cur++
					break
				}
				r.cur++
			}
			fmt.Println("out of comment at:", r.lineNo, string(r.lineNo))
			continue
		}

		// is word?
		if unicode.IsLetter(c) {
			n := 0
			for ; unicode.IsLetter(r.peekN(n)) ||
				unicode.IsDigit(r.peekN(n)); n++ {
			}
			if w := string(r.text[r.cur : r.cur+n]); w == "doğru" || w == "yanlış" {
				r.consume(n, "bool")
				continue
			}
			r.consume(n, "word")
			continue
		}

		// is no?
		if unicode.IsDigit(c) {
			// read till the end
			n := 0
			for ; unicode.IsDigit(r.peekN(n)); n++ {
			}
			r.consume(n, "number")
			continue
		}

		// is punct?
		if c == ':' {
			r.consume(1, ":")
			continue
		}

		if c == '<' || c == '>' {
			if r.peek() == '=' {
				r.consume(2, "op")
				continue
			}
			r.consume(1, "op")
			continue
		}

		// if c == 'VE' || c == 'VEYA'

		// ar for arithmeic
		if c == '+' || c == '-' || c == '*' || (c == '/') {
			r.consume(1, "ar")
			continue
		}

		// bracket
		if c == '(' || c == ')' {
			r.consume(1, "br")
			continue
		}

		if c == '{' || c == '}' {
			r.consume(1, "par")
			continue
		}

		if c == '\n' {
			r.consume(1, "nl")
			r.lineNo++
			continue
		}

		if c == '=' {
			r.consume(1, "eq")
			continue
		}

		// in this case look prev
		// Belirtme veya Yükleme Hali: Belirtme hali ekleri -ı, -i, -u, -ü
		// "y" kaynastirma
		if c == '\'' {
			if b := r.backstep(); unicode.IsLetter(b) || unicode.IsDigit(b) {
				// read until next whitespace
				n := 0
				for ; r.peekN(n) != ' '; n++ {
				}
				p := "[y,n]?[i,ı,ü,u]"
				ek := r.text[r.cur : r.cur+n]
				match, _ := regexp.MatchString(p, string(ek))
				if !match {
					log.Fatal("wrong postfix at:", r.cur)
				}
				r.consume(n, "ek")
			}
		}

		if c == '.' {
			r.consume(1, "dot")
			continue
		}

		r.cur++

	}

}

func (r *reader) now() rune {
	return r.text[r.cur]
}

func (r *reader) peek() rune {
	return r.text[r.cur+1]
}

func (r *reader) backstep() rune {
	return r.text[r.cur-1]
}

func (r *reader) peekN(n int) rune {
	return r.text[r.cur+n]
}

func (r *reader) consume(step int, kind string) {
	t := token{
		val:  string(r.text[r.cur : r.cur+step]),
		kind: kind,
		pos:  r.cur,
		line: r.lineNo,
	}
	r.cur += step
	r.tokenList = append(r.tokenList, t)
}
