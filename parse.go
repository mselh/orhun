package main

import (
	"log"
)

type node struct {
	val  string
	kind string
	line int

	left  *node
	right *node

	// only used for block kind node
	exprs []*node

	// only used for "fnParams" node
	fnParams []*node

	// only used with if statement
	ifNode *node
}

func (n *node) Print() {
	nprint(n, "")
}

func nprint(n *node, space string) {
	if n == nil {
		return
	}

	myPrintln(space, "---node---")
	myPrintln(space, n.kind)
	myPrintln(space, n.val)
	if len(n.fnParams) > 0 {
		myPrint(space, "fn params:")
		for i := range n.fnParams {
			myPrint(n.fnParams[i].val, " ")
		}
		myPrintln()
	}
	if len(n.exprs) > 0 {
		myPrintln(space, "expr:")
		for i := range n.exprs {
			nprint(n.exprs[i], space+" ")
		}
		myPrintln()
	}
	oldspace := space
	space = space + " "
	nprint(n.left, space)
	nprint(n.right, space)
	myPrintln(oldspace, "---end---")
}

type parser struct {
	tokenList []token
	cur       int
	root      *node
}

func (p *parser) parseAll() {

	// store global new vars in the expressions[] node of the root
	p.root.exprs = make([]*node, 0)

	for p.cur < len(p.tokenList) {
		t := p.tokenList[p.cur]
		myPrintln("parsing tok:", t)

		if t.val == "giriş" {
			entry := p.parseEntry()
			entry.Print()
			// below may change
			p.root.left = entry
			continue
		}

		if t.kind == "nl" {
			p.cur++
			continue
		}

		// parse global "yeni" as expressions
		if t.val == "yeni" {
			nvar := p.parseNew()
			if nvar == nil {
				log.Fatalln("err while parsing new", t.line)
			}
			p.root.exprs = append(p.root.exprs, nvar)
			continue
		}

		myPrintln("cannot parse:", t)
		return
	}
}

func (p *parser) parseEntry() *node {
	// skip giris
	p.cur++
	return p.parseBlock("entry")
}

// expressions
// starts with yeni
// starts with an already defined identifier
// can be a module name, or a variable

// starts with eger
// starts with yinele
func (p *parser) parseExpr() *node {

	// skip starting new lines

	// look for kw
	if p.now().val == "yeni" {
		return p.parseNew()
	}

	// if
	// eger [bool.value] dogruysa:
	// block
	// .?
	//}
	// left note if true block
	// r is else
	//
	// for loop
	// eger bool.value dogruysa yinele:
	// block
	// .
	if p.now().val == "eğer" {
		n := new(node)
		n.kind = "if"
		p.cur++
		n.ifNode = p.parseVal()
		p.cur++ // skip dogruysa

		// check for while
		if p.now().val == "yinele" {
			n.kind = "while"
			p.cur++ //skip yinele
			n.left = p.parseBlock("while")
			return n
		}

		n.left = p.parseBlock("ifr")
		if p.now().kind == "nl" {
			p.cur++
		}

		if p.now().val == "değilse" {
			p.cur++
			n.right = p.parseBlock("ifw")
		}

		return n
	}

	// syntax
	// for, variable = new value , act as assing
	// for, [v'[postfix]]? [inscope function]
	// evaluate
	if p.now().kind == "word" {
		if p.peekN(1).val == "=" {
			return p.parseAssign()
		}
		return p.parseFnCall()
	}

	// if nothing is matched give an error
	log.Fatal("expression is unkown at:", p.cur, p.now().val, "val", []rune(p.now().val))
	return nil
}

// the subexpr parsers returns node, one hould at them to root
// parses new variables and new functions
func (p *parser) parseNew() *node {
	if p.tokenList[p.cur+3].val != "=" {
		log.Fatal("bad new variable decl. @", p.cur)
		return nil
	}

	varname := p.tokenList[p.cur+2]
	typename := p.tokenList[p.cur+1]

	p.cur = p.cur + 4
	var value *node
	if typename.kind != "fn" {
		value = p.parseVal()
	} else {
		value = p.parseFnDef()
	}

	n := new(node)
	n.kind = "new"
	n.val = varname.val
	// bence deger ve tipi tutulmali,
	// sonra tip ve atamayi karsilastirmak gerekebilir
	// nizami olmasini kontrol etmek isteyebiliriz.

	n.left = new(node)
	n.left.kind = "typeName"
	n.left.val = typename.val

	n.right = value
	// while traversing tree, one should eval this
	return n
}

// what is val?
// either be int or bool for now
// single constant number/integer
// or math arithmetic
// SAYI: SAYI
// SAYI: SAYI +-/* SAYI
// ONERME: ONERME (VE VEYA ONERME)*
// ONERME: [DOGRU | YANLIS]
// ONERME:  SAYI < > <= => SAYI
func (p *parser) parseVal() *node {

	// parantheses
	if p.now().val == "(" {
		n := new(node)
		n.kind = "PAR"
		p.cur++
		n.left = p.parseVal()
		if p.now().val != ")" {
			log.Fatal("unclosed paranthesis", p.now())
		}
		p.cur++
		return n
	}

	if p.peekN(0).kind == "word" {
		n := new(node)
		n.kind = "VAR"
		n.val = p.now().val
		p.cur++ // skip word
		//myPrintln("next:", p.now())

		// X VE|VEYA Y
		if p.now().kind == "op" {
			root := new(node)
			root.kind = "OP"
			root.val = p.now().val

			// skip op
			p.cur++

			right := p.parseVal()
			root.right = right
			root.left = n
			return root
		}

		if p.now().kind == "rel" {
			root := new(node)
			root.kind = "REL"
			root.val = p.now().val

			// skip op
			p.cur++

			right := p.parseVal()
			root.right = right
			root.left = n
			return root
		}

		// X [[+/*-] Y ]*
		if p.now().kind == "ar" {
			root := new(node)
			root.kind = "AR"
			root.val = p.now().val

			// skip op
			p.cur++

			right := p.parseVal()
			root.right = right
			root.left = n
			return root
		}

		return n
	}

	if p.now().kind == "number" {
		n := new(node)
		n.kind = "NUM"
		n.line = p.now().line
		n.val = p.now().val
		p.cur++
		return n
	}
	if p.now().kind == "bool" {
		n := new(node)
		n.kind = "BOOL"
		n.val = p.now().val
		p.cur++
		return n
	}

	log.Fatal("couldnt parse at: ", p.cur, p.now())
	return nil
}

func (p *parser) parseAssign() *node {

	// check
	if p.peekN(1).val != "=" {
		log.Fatal("not a valid assignment:", p.now())
	}

	n := new(node)
	n.kind = "assign"

	n.left = new(node)
	n.left.kind = "var"
	n.left.val = p.now().val

	p.cur += 2
	n.right = p.parseVal()

	return n
}

// parsing a function
var fnParam struct {
	name     string
	typename string
}

// parses fn def
func (p *parser) parseFnDef() *node {
	n := new(node)
	n.kind = "fn"

	// look for ()
	if p.now().val != "(" {
		log.Fatalln("bad fn definition at:", p.now().line, p.now())
	}
	p.cur++

	// parse input params ([name type[,name type]*]?)
	start := p.now()
	for p.now().val != ")" {
		if p.cur >= len(p.tokenList) {
			log.Fatalln("bad fn input defn. no ) after", start.line, start)
		}

	}

	return n
}

// word'ek fn
// [word'ek ve]* word'ek fn
// fn
// ilerde buraya bool ve int konmasini dusunebiliriz.
func (p *parser) parseFnCall() *node {
	n := new(node)
	n.kind = "fnCall"

	n.right = new(node)
	n.right.val = "params"
	n.right.fnParams = make([]*node, 0)

	// get functions pos
	fnStart := 0
	for p.peekN(fnStart).kind == "word" &&
		p.peekN(fnStart+1).kind == "ek" {
		// make ek optional, feedback from an old friend

		// add fn params
		pn := new(node)
		pn.kind = "var"
		pn.val = p.peekN(fnStart).val

		n.right.fnParams = append(n.right.fnParams, pn)
		if p.peekN(fnStart+2).val == "ve" {
			fnStart += 3
		} else {
			fnStart += 2
		}
	}

	n.left = new(node)
	n.left.kind = "fn"
	n.left.val = p.peekN(fnStart).val

	p.cur += fnStart + 1
	p.cur++
	return n
}

func (p *parser) parseBlock(blockType string) *node {

	if p.now().val != ":" {
		log.Fatal("cant parse block start at", p.now())
	}

	p.cur++
	b := new(node) // new block node
	b.kind = "block"
	b.val = blockType
	b.exprs = make([]*node, 0)

	// read all
	for p.now().val != "." {
		if p.now().kind == "nl" {
			p.cur++
			continue
		}
		//myPrintln("parsing expr starting at:", p.now())
		n := p.parseExpr()
		//n.Print()
		b.exprs = append(b.exprs, n)
	}
	p.cur++ // skip last dot
	return b
}

// helpers
func (p *parser) now() token {
	return p.tokenList[p.cur]
}

func (p *parser) peekN(n int) token {
	return p.tokenList[p.cur+n]
}
