package main

import (
	"fmt"
	"log"
	"strconv"
)

type program struct {
	rootNode  *node // root ast node
	rootScope *scope
}

type scope struct {
	parent    *scope
	localVars map[string]*val
}

func (s *scope) String() string {
	str := ""
	for k, v := range s.localVars {
		str += fmt.Sprintln(k, *v)
	}
	return str
}

type Obj struct {
	name string
	// both acts as a val and kind
	keys map[string]*val
}

type fn struct {
	// if inherited from golang add a reference
	GoRef func([]*val) []val
	// if orhun defined function add a reference to node
	// so that the fnDef node is executed with the params
	isParamsParametric bool
	signature          []val // use this to verify params
	defNode            *node // node to be executed
}

func (f *fn) exec(params []*val) ([]val, error) {

	// check signature
	// check if params is in same kind order as signature
	if !f.isParamsParametric {
		if len(params) != len(f.signature) {
			return nil, fmt.Errorf("bad func signature")
		}
	} else {
		if len(params) == 0 {
			return nil, fmt.Errorf("at least a signle param req.")
		}
	}

	for i := range f.signature {
		sv := f.signature[i]
		if sv.typeName != params[i].typeName {
			return nil, fmt.Errorf("bad func parameter")
		}
	}
	// check done

	// for now skip orhun defined functions
	if f.GoRef != nil {
		return f.GoRef(params), nil
	}
	return nil, nil
}

type val struct {
	// int,string,object
	typeName string
	// if int
	intval int
	// if string
	strval string
	// if bool
	boolVal bool
	// if ozne, meaning object, object reference
	objVal Obj
	// if val of type fn, it means function
	funcVal fn
}

func (p *program) walk() {
	// start by reading deneme.orhun
	// module requires multiple files
	// it is a milestone for future

	// defs outside the program is not supported now,
	// yeni keyword might be a good starting point
	// ie. yeni islev
	// yeni sabit
	// yeni tamsayi ...

	// add builtin variables to root scope
	addBuiltins(p.rootScope)

	// start from entry
	entryNode := p.rootNode.left
	entryScope := new(scope)
	entryScope.parent = p.rootScope
	entryScope.localVars = make(map[string]*val)
	for i := range entryNode.exprs {
		e := entryNode.exprs[i]
		// forgets children scopes after exec
		def, val := exec(e, entryScope)
		if def != "" {
			entryScope.localVars[def] = val
			log.Println("added new def to scope:", entryScope.localVars)
		}

	}

}

// execs a node
// if it s a new node, returns that new node
// else only executes it.
func exec(e *node, parentScope *scope) (def string, v *val) {
	if e.kind == "new" {
		name := e.val // expression val is the name
		v := new(val)
		if e.left.val == "tamsayı" {
			v.typeName = "int"
			v.intval = evalIntValue(e.right)
		}

		return name, v
	}
	if e.kind == "assign" {
		// search for the var in local scope
		// ifndef, go for upper scope
		v := searchVar(parentScope, e.left.val)
		if v == nil {
			log.Fatal("cant find variable in the scope, can't assign, ", e.line, e)
		}

		if v.typeName == "int" {
			v.intval = evalIntValue(e.right)
		}

		fmt.Println("reassigned var:", e.left.val)
		fmt.Println(parentScope)
		return "", nil
	}

	if e.kind == "fnCall" {
		// for some functions we will use golang fns
		// for some it will execute the function that is defined
		// for start, start with golang and builtin
		fnName := e.left.val
		v := searchVar(parentScope, fnName)
		if v == nil {
			log.Fatal("cant find fn variable in the scope,", fnName, e.line, e)
		}
		if v.typeName != "fn" {
			log.Fatal("not a function", e.line, e)
		}
		// start getting fnParams from right node
		fnParams := make([]*val, 0)
		for i := range e.right.fnParams {
			pn := e.right.fnParams[i]
			paramVar := searchVar(parentScope, pn.val)
			if paramVar == nil {
				log.Fatal("can t find fn param var in the scope", fnName, pn.val, e.line, e)
			}
			fnParams = append(fnParams, paramVar)
		}
		_, err := v.funcVal.exec(fnParams)
		if err != nil {
			log.Fatal("err while exec in g fn,", err, fnName, e.line, e)
		}
		return "", nil

	}

	log.Fatal("can't exec current node:", e.line, e)
	return "", nil
}

func evalIntValue(n *node) int {
	//walks and evaluates a valude node, who is int
	if n.kind == "NUM" {
		i, err := strconv.Atoi(n.val)
		if err != nil {
			log.Fatal("couldn't evaluate num at", n.line, err)
		}
		return i
	}

	log.Fatal("can't parse node at", n.line, n)
	return 0
}

// helpers
func newProgram(root *node) *program {
	rootScope := new(scope)
	rootScope.localVars = make(map[string]*val)

	return &program{
		rootNode:  root,
		rootScope: rootScope,
	}
}

func searchVar(cur *scope, key string) *val {
	v, found := cur.localVars[key]
	for !found {
		cur = cur.parent
		if cur == nil {
			return nil
		}
		v, found = cur.localVars[key]
	}
	return v
}

func addBuiltins(s *scope) {
	printFn := new(fn)
	printFn.GoRef = func(v []*val) []val {
		fmt.Println("WRITE CALLED")
		for i := range v {
			fmt.Print(v[i].intval, " ")
		}
		fmt.Println("\nWRITE FINISHED")
		return nil
	}
	printFn.isParamsParametric = true
	printFn.signature = []val{{typeName: "int"}}

	printFnVal := new(val)
	printFnVal.typeName = "fn"
	printFnVal.funcVal = *printFn
	s.localVars["yazdır"] = printFnVal

}
