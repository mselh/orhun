package main

import (
	"fmt"
	"strconv"
)

type Node struct {
	lhs *Node // left hand side
	rhs *Node // right hand side

	// Block, used if kind == ND_BLOCK
	body *Node

	//
	kind string

	//
	variable *Obj // used if kind == ND_VAR

	val int // used if kind = VAL
}

// local variable
type Obj struct {
	name   string
	ty     string
	offset int // offset from rbp, think this as addr but for local stack
}

type Function struct {
	body   []*Node
	locals []*Obj
}

func (f *Function) newLVar(name string, ty string) *Obj {
	v := new(Obj)
	v.name = name
	v.ty = ty
	f.locals = append(f.locals, v)
	return v
}

func declaration(tokenList []token, tok token, f *Function, start int) (declNode *Node, skip int) {
	head := &Node{}

	tokenList = tokenList[start:]
	i := 0
	for ; tokenList[i].val != "\n"; i++ {
		if i == 1 {
		}

		if i == 2 {
			// variable type
			v := f.newLVar(tokenList[1].val, tokenList[2].val)
			lhs := new(Node)
			lhs.variable = v
			head.lhs = lhs
		}

		if i == 3 {
			if tokenList[3].val != "=" {
				fmt.Println("tok: ", tok.val)
				println(tokenList)
				panic("=?")
			}
		}

		if i == 4 {
			// value
			// head
			rhs := new(Node)
			val, err := strconv.Atoi(tokenList[4].val)
			if err != nil {
				panic(err)
			}
			rhs.val = val
			head.rhs = rhs
		}

	}
	fmt.Println("parsed, ", tokenList[:i])
	skip = i

	return head, skip
}

// only look for num for now
func expression(tl []token) (node *Node, skip int) {
	return assign(tl)
}

// assign = equality
func assign(tl []token) (node *Node, skip int) {
	node := equality(tl)
	return node
}

func add(tl []token) *Node {
	node := mul(tl)

	i := 0
	for ; tl[i].val != "\n"; i++ {
		if tl[i].val == "+" {
			node = NewAdd(node, mul(tl))
			continue
		}

		return node
	}
}

// mul = unary ("*" unary | "/" unary)*
// for now mul = unary
func mul(tl []token) *Node {
	node := unary(tl)
	return node
}

func unary(tl []token) *Node {
	return primary(tl)
}

// find local var by name
func (f *Function) findVar(string name) *Obj {
	for _, v := range f.locals {
		if v.name == name {
			return v
		}
	}
	return nil
}

// primary
// needs to access local variables of the function
func primary(tl []token, f *Function) (primary *Node, skip int) {
	if tl[1].val == "IDENTIFIER" {
		// Variable
		v := findVar(tl[1].val)
		if v == nil {
			panic(tok, "undefined variable")
		}
		skip = 1
		n := new(Node)
		n.kind = "VARIABLE"
		n.variable = v
		return n, skip
	}

	if tok.kind == "NUMBER" {
		node := NewNum(tok.val, tok)
		*rest = tok.Next
		return node
	}

	panic("expected an expression")
}

func NewAdd(lhs, rhs *Node, tl []token) (binaryNode *Node) {
	if lhs.kind == "INTEGER" && rhs.kind == "INTEGER" {
		n := new(Node)
		n.kind = "ADD"
		n.lhs = lhs
		n.rhs = rhs
		return n
	}

	panic("bad add expr")
}

// equality = relational
func equality(tl []token) (node *Node, skip int) {
	return relational()
}

func relational(tl []token) (node *Node, skip int) {
	return add(tl)
}

func compoundStmt(tl []token, f *Function) []*Node {
	stmtList := make([]*Node, 0)
	skip := 0
	for i, tok := range tl {
		if skip > 0 {
			continue
		}
		if tok.val == "değişken" {
			var decl *Node
			ret.kind = "DECLARATION"
			decl, skip = declaration(tl, tok, f, i)
			stmtList = append(stmtList, decl)
		} else if tok.val == "döndür" {
			var ret *Node
			ret.kind = "RETURN"
			ret.lhs, skip = expression(tl, i+1)
			stmtList = append(stmtList, ret)
		}
	}
	return stmtList
}
