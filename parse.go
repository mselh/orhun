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

type Obj struct {
	name string
	ty   string
	next *Obj
}

type Function struct {
	body   []*Node
	locals *Obj
}

func (f *Function) newLVar(name string, ty string) *Obj {
	v := new(Obj)
	v.name = name
	v.ty = ty
	v.next = f.locals
	f.locals = v
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
func expression(tl []token) {

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
