package main

import (
	"log"
	"regexp"
	"unicode"
)

type token struct {
	kind string
	val  string
	pos  int
}

type reader struct {
	text      []rune
	cur       int
	tokenList []token
}

func (r *reader) tokenize() {
	// dont token whitespaces
	for r.cur < len(r.text) {
		c := r.now()

		// is word?
		if unicode.IsLetter(c) {
			n := 0
			for ; unicode.IsLetter(r.peekN(n)) ||
				unicode.IsDigit(r.peekN(n)); n++ {
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

		// ar for arithmeic
		if c == '+' || c == '-' || c == '*' || c == '/' {
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
	}
	r.cur += step
	r.tokenList = append(r.tokenList, t)
}
