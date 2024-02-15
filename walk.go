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
	funcVal *fn
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

	// add global declaratins to root scope
	for i := range p.rootNode.exprs {
		n := p.rootNode.exprs[i]
		if n.kind != "new" {
			log.Fatalln("unexpected exprs,", n.line, "only 'yeni' is allowed")
		}
		key, v := exec(n, p.rootScope)
		if searchVar(p.rootScope, key) != nil {
			log.Fatalln("already defined key,", key)
		}
		p.rootScope.localVars[key] = v
	} // TODO: evaluate them later if right side includes a variable
	// you also need to detect recursive declarations etc. too much work for now
	myPrintln("GLOBAL SCOPE:", p.rootScope)

	// start from entry
	entryNode := p.rootNode.left
	if p.rootNode.left == nil {
		log.Fatalln("no entry node giris defined!")
	}
	execBlockNode(entryNode, p.rootScope)

}

// execs a node
// if it s a new node, returns that new node
// else only executes it.
func exec(e *node, parentScope *scope) (def string, v *val) {
	if e.kind == "new" {
		name := e.val // expression val is the name
		v := new(val)
		v.typeName = e.left.val

		// if FN Definition walk the node and return the fndef val
		// this is to add new fn definitions
		if e.left.val == "işlev" {
			if e.right.kind != "fnDefn" {
				log.Fatal("waiting for fn def.")
			}
			f := new(fn)
			f.signature = make([]val, 0)
			for i := range e.right.fnSignature {
				fns := e.right.fnSignature[i]
				v := val{
					typeName: fns.typename,
				}
				f.signature = append(f.signature, v)
			}
			f.defNode = e.right
			v.typeName = "fn"
			v.funcVal = f
		}

		// FN CALL
		if e.right.kind == "fnCall" {
			vals := execFnNode(e.right, parentScope)
			// below is a placeholder
			// also check the left side
			// a function call should return a single value now
			// whether it is int, bool or a new function
			// when implemented, also structs
			if len(vals) > 0 {
				v.intval = vals[0].intval
			}
			fmt.Println("fn call vals:", vals)
			return name, v
		}
		// FN CALL

		// initial values are evaluated according to the type
		if e.left.val == "tamsayı" {
			v.intval = evalIntValue(e.right, parentScope)
		}
		if e.left.val == "önerme" {
			v.boolVal = evalBoolValue(e.right, parentScope)
		}
		myPrintln("new var:", e.left.val)
		return name, v
	}
	if e.kind == "assign" {
		// search for the var in local scope
		// ifndef, go for upper scope
		v := searchVar(parentScope, e.left.val)
		if v == nil {
			log.Fatal("cant find variable in the scope, can't assign, ", e.line, e)
		}

		if v.typeName == "tamsayı" {
			v.intval = evalIntValue(e.right, parentScope)
		}
		if v.typeName == "önerme" {
			v.boolVal = evalBoolValue(e.right, parentScope)
			myPrintln("bool:", v.boolVal)
		}

		myPrintln("reassigned var:", e.left.val)
		myPrintln(parentScope)
		return "", nil
	}

	if e.kind == "fnCall" {
		execFnNode(e, parentScope)
		return "", nil

	}

	if e.kind == "if" {
		execIfNode(e, parentScope)
		myPrintln(parentScope)
		return "", nil
	}

	if e.kind == "while" {
		myPrintln("entered while")
		execWhileNode(e, parentScope)
		myPrintln(parentScope)
		return "", nil
	}

	log.Fatal("can't exec current node:", e.line, e)
	return "", nil
}

func evalIntValue(n *node, s *scope) int {
	//walks and evaluates a valude node, who is int
	if n.kind == "NUM" {
		i, err := strconv.Atoi(n.val)
		if err != nil {
			log.Fatal("couldn't evaluate num at", n.line, err)
		}
		return i
	}

	if n.kind == "VAR" {
		myPrintln("evaling VAR node as int,", n.line, n)
		v := searchVar(s, n.val)
		return v.intval
	}

	if n.kind == "PAR" {
		// evaluate the pharantesis then return
		return evalIntValue(n.left, s)
	}

	if n.kind == "AR" {

		leftval := evalIntValue(n.left, s)
		rightval := evalIntValue(n.right, s)
		if n.val == "+" {
			return leftval + rightval
		}
		if n.val == "*" {
			return leftval * rightval
		}
		if n.val == "/" {
			return leftval / rightval
		}
		if n.val == "-" {
			return leftval - rightval
		}
	}

	log.Fatal("can't walk node at", n.line, n)
	return 0
}

func evalBoolValue(n *node, s *scope) bool {
	//walks and evaluates a value node, who is bool
	if n.kind == "BOOL" {
		return n.val == "doğru"
	}

	if n.kind == "PAR" {
		// evaluate the pharantesis then return
		return evalBoolValue(n.left, s)
	}

	if n.kind == "VAR" {
		myPrintln("evaling VAR node as bool,", n.line, n)
		v := searchVar(s, n.val)
		return v.boolVal
	}

	if n.kind == "OP" {

		//if n.left.kind == "BOOL" && n.right.kind == "BOOL" {
		leftval := evalBoolValue(n.left, s)
		rightval := evalBoolValue(n.right, s)
		if n.val == "ve" {
			return leftval && rightval
		}
		if n.val == "veya" {
			return leftval || rightval
		}
		//}

	}

	if n.kind == "REL" {
		leftval := evalIntValue(n.left, s)
		rightval := evalIntValue(n.right, s)
		if n.val == "<" {
			return leftval < rightval
		}
		if n.val == "<=" {
			return leftval <= rightval
		}
		if n.val == ">" {
			return leftval > rightval
		}
		if n.val == ">=" {
			return leftval >= rightval
		}
	}

	n.Print()
	log.Fatalln("cant evaluate bool value of node:", n.line)
	return false
}

func execIfNode(n *node, s *scope) {
	// evaluate if expression
	isTrue := evalBoolValue(n.ifNode, s)
	blockNode := n.right
	if isTrue {
		blockNode = n.left
	}

	execBlockNode(blockNode, s)
}

func execWhileNode(n *node, s *scope) {
	//eval if in every loop
	for {
		if !evalBoolValue(n.ifNode, s) {
			break
		}

		execBlockNode(n.left, s)
	}
}

func execFnNode(e *node, s *scope) []val {
	// for some functions we will use golang fns
	// for some it will execute the function that is defined
	// for start, start with golang and builtin
	fnName := e.left.val
	v := searchVar(s, fnName)
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
		paramVar := searchVar(s, pn.val)
		if paramVar == nil {
			log.Fatal("can t find fn param var in the scope", fnName, pn.val, e.line, e)
		}
		fnParams = append(fnParams, paramVar)
	}
	vars, err := v.funcVal.exec(fnParams)
	if err != nil {
		log.Fatal("err while execing fn,", err, fnName, e.line, e)
	}
	return vars
}

// add return val, dondur function
func execBlockNode(n *node, parentScope *scope) *val {
	blockScope := new(scope)
	blockScope.parent = parentScope
	blockScope.localVars = make(map[string]*val)
	for i := range n.exprs {
		e := n.exprs[i]
		// forgets children scopes after exec
		def, val := exec(e, blockScope)
		if def != "" {
			blockScope.localVars[def] = val
			myPrintln("added new def to scope :", def, *val, blockScope.localVars)
		}
	}
	return nil
}

// helpers
func newProgram(root *node) *program {
	rootScope := new(scope)
	rootScope.localVars = make(map[string]*val)
	rootScope.parent = nil

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
		myPrintln("WRITE CALLED")
		for i := range v {
			fmt.Print(v[i].intval, " ")
		}
		fmt.Println()
		myPrintln("\nWRITE FINISHED")
		return nil
	}
	printFn.isParamsParametric = true
	printFn.signature = []val{{typeName: "tamsayı"}}

	printFnVal := new(val)
	printFnVal.typeName = "fn"
	printFnVal.funcVal = printFn
	s.localVars["yazdır"] = printFnVal

}
