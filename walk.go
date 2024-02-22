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
	GoRef     func([]*val) []*val
	rootScope *scope
	// if orhun defined function add a reference to node
	// so that the fnDef node is executed with the params
	isParamsParametric bool
	signature          []val // use this to verify params
	defNode            *node // node to be executed
}

func (f *fn) exec(params []*val) ([]*val, error) {

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
		if sv.typeName != "" && sv.typeName != params[i].typeName {
			return nil, fmt.Errorf("bad func parameter")
		}
	}
	// check done

	if f.GoRef != nil {
		return f.GoRef(params), nil
	}

	// exec fn block

	// add parameters to local scope
	fnScope := new(scope)
	fnScope.localVars = make(map[string]*val)
	fnScope.parent = f.rootScope
	for i := range params {
		p := params[i]
		myPrintln("param", *p)
		key := f.defNode.fnSignature[i]
		fnScope.localVars[key.name] = p
	}
	myPrint("FN SCOPE:")
	myPrintln(fnScope)

	returnTypes := []string{f.defNode.returnTypename}
	// walk block
	vals := execBlockNode(f.defNode.left, fnScope, returnTypes)

	// look for returned values

	return vals, nil
}

type val struct {
	// int,string,struct,fn
	typeName string
	// if int
	intval int
	// if string
	strval string
	// if bool
	boolVal bool
	// if yapi, meaning struct
	objVal *Obj
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

		if n.kind != "new" && n.kind != "yapı" {
			log.Fatalln("unexpected exprs,", n, n.line, "only 'yeni' is allowed")
		}

		if n.kind == "new" {

			key, v := exec(n, p.rootScope, make([]string, 0))
			if searchVar(p.rootScope, key) != nil {
				log.Fatalln("already defined key,", key)
			}
			p.rootScope.localVars[key] = v
		}

		if n.kind == "yapı" {
			key, v := exec(n, p.rootScope, make([]string, 0))
			if searchVar(p.rootScope, key) != nil {
				log.Fatalln("already defined key,", key)
			}
			p.rootScope.localVars[key] = v
			myPrintln("new yapi with fields", v.objVal.keys)
		}

	} // TODO: evaluate them later if right side includes a variable
	// you also need to detect recursive declarations etc. too much work for now
	// for now global declatrations are top-down

	myPrintln("GLOBAL SCOPE:", p.rootScope)

	// start from entry
	entryNode := p.rootNode.left
	if p.rootNode.left == nil {
		log.Fatalln("no entry node giris defined!")
	}
	execBlockNode(entryNode, p.rootScope, make([]string, 0))

}

// execs a node
// if it s a new node, returns that new node
// else only executes it.
func exec(e *node, parentScope *scope, retval []string) (def string, v *val) {

	if e.kind == "yapı" {
		myPrintln("ypai val:", e.val)
		name := e.val
		v := new(val)
		v.typeName = "yapıDefn"
		v.objVal = new(Obj)
		v.objVal.name = name
		v.objVal.keys = make(map[string]*val)
		for i := range e.fnSignature {
			defn := e.fnSignature[i]
			val := new(val)
			val.typeName = defn.typename
			v.objVal.keys[defn.name] = val
		}
		return name, v
	}

	// below is := style declerations
	// type assignment is done by looking at right side
	// fn defns are not implemented for := operator
	if e.kind == "newAssign" {
		name := e.val
		if e.right.kind == "fnCall" {
			vals := execFnNode(e.right, parentScope)
			return name, vals[0]
		} else {
			v := searchVar(parentScope, e.right.val)
			if v != nil && v.typeName == "fn" {
				fnNode := new(node)
				fnNode.left = new(node)
				fnNode.left.val = e.right.val
				fnNode.right = v.funcVal.defNode
				vals := execFnNode(fnNode, parentScope)
				return name, vals[0]
			}
		}

		if e.right.kind == "OP" {
			v := new(val)
			v.typeName = "önerme"
			v.boolVal = evalBoolValue(e.right, parentScope)
			return name, v
		}

		if e.right.kind == "VAR" {
			v1 := searchVar(parentScope, e.right.val)
			if v == nil {
				log.Fatalln("undefined VAR", e.right, e.line)
			}
			v2 := new(val)
			v2.typeName = v1.typeName
			v2.intval = v1.intval
			v2.boolVal = v1.boolVal
			v2.objVal = v1.objVal
			return name, v2
		}

		if e.right.kind == "REL" {
			v := new(val)
			v.typeName = "önerme"
			v.boolVal = evalBoolValue(e.right, parentScope)
			return name, v
		}

		if e.right.kind == "FIELD" {
			// get defined object
			log.Fatalln("not yet implemented")
			// then get the field
			// then return the fields value
		}

		if e.right.kind == "AR" {
			v := new(val)
			v.typeName = "tamsayı"
			v.intval = evalIntValue(e.right, parentScope)
			return name, v
		}

		if e.right.kind == "NUM" {
			v := new(val)
			v.typeName = "tamsayı"
			v.intval = evalIntValue(e.right, parentScope)
			return name, v
		}

		if e.right.kind == "BOOL" {
			v := new(val)
			v.typeName = "önerme"
			v.boolVal = evalBoolValue(e.right, parentScope)
			return name, v
		}

		if e.right.kind == "STRING" {
			v := new(val)
			v.typeName = "metin"
			v.strval = e.right.val
			return name, v
		}

	}

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
			f.rootScope = parentScope
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

		// not initialized var with no right side
		if e.right == nil {
			// search for the key
			if searchVar(parentScope, name) != nil {
				log.Fatalln("already defined new:", name, e.line)
			}
			// check the typename and return with empty vals:
			zeroed := new(val)
			if v.typeName == "tamsayı" || v.typeName == "önerme" {
				zeroed.typeName = v.typeName
			}
			// else check if it matched a defined struct type
			defnVal := searchVar(parentScope, v.typeName)
			if defnVal == nil {
				log.Fatalln("undefined typename:", v.typeName, e.line)
			}
			if defnVal.typeName != "yapıDefn" {
				log.Fatalln("not a defined struct typename:", v.typeName, e.line)
			}
			// create a new instance from defn
			zeroed.typeName = "yapı"
			zeroed.objVal = new(Obj)
			zeroed.objVal.name = defnVal.objVal.name
			zeroed.objVal.keys = make(map[string]*val)
			for k, v := range defnVal.objVal.keys {
				zeroed.objVal.keys[k] = new(val)
				zeroed.objVal.keys[k].typeName = v.typeName
			}
			return name, zeroed
		}

		// FN CALL
		if e.right.kind == "fnCall" {
			vals := execFnNode(e.right, parentScope)
			// below is a placeholder
			// also check the left side
			// a function call should return a single value now
			// whether it is int, bool or a new function
			// when implemented, also structs

			// the below should change to multiple vals
			// in future
			if len(vals) > 0 {
				v.intval = vals[0].intval
			}
			return name, v
		} else {
			// check if right is defined in scope as fn
			v := searchVar(parentScope, e.right.val)
			if v != nil && v.typeName == "fn" {
				// new fn node is required to use execfnNode
				fnNode := new(node)
				fnNode.left = new(node)
				fnNode.left.val = e.right.val
				fnNode.right = v.funcVal.defNode
				//v.funcVal.defNode.Print()
				vals := execFnNode(fnNode, parentScope)
				return name, vals[0]
			}
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

	if e.kind == "fieldAssign" {
		// first find the struct var
		v := searchVar(parentScope, e.val)
		if v == nil {
			log.Fatalln("cant assign to not defined struct", e.val)
		}

		fieldName := e.left.left.val

		// assign according to the type
		fieldType := v.objVal.keys[fieldName].typeName
		if fieldType == "tamsayı" {
			v.objVal.keys[fieldName].intval = evalIntValue(e.left.right, parentScope)
		}
		if fieldType == "önerme" {
			v.objVal.keys[fieldName].boolVal = evalBoolValue(e.left.right, parentScope)
		}
		if fieldType == "metin" {
			v.objVal.keys[fieldName].strval = e.left.right.val
		}
		return "", nil

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
		execIfNode(e, parentScope, retval)
		myPrintln(parentScope)
		return "", nil
	}

	if e.kind == "while" {
		myPrintln("entered while")
		execWhileNode(e, parentScope, retval)
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

	if n.kind == "FIELD" {
		myPrintln("evaling FIELD node as int,", n.line, n)
		v := searchVar(s, n.val)
		if v == nil {
			log.Fatalln("cant find var", n.val, n.line)
		}
		fieldname := n.left.val
		fieldval, found := v.objVal.keys[fieldname]
		if !found {
			log.Fatalln("cant find field at struct", n.val, n.line)
		}
		if fieldval.typeName != "tamsayı" {
			log.Fatalln("field is not an integer:", fieldname, n.line)
		}
		intval := fieldval.intval

		return intval
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

func evalReturnVal(n *node, s *scope, returnTypes []string) []*val {
	// node here is an fn node
	if len(returnTypes) == 0 {
		return make([]*val, 0)
	}

	retVals := make([]*val, 0)
	for i := range n.right.fnParams {
		n := n.right.fnParams[i]
		n.Print()
		v := new(val)

		v.typeName = n.val
		if returnTypes[i] == "önerme" {
			v.boolVal = evalBoolValue(n, s)
		} else if returnTypes[i] == "tamsayı" {
			v.intval = evalIntValue(n, s)
		} else {
			log.Fatalln("ret val is not int or bool", returnTypes[i], n.line, n)
		}
		v.typeName = returnTypes[i]
		retVals = append(retVals, v)
	}

	return retVals
}

func execIfNode(n *node, s *scope, retVal []string) {
	// evaluate if expression
	isTrue := evalBoolValue(n.ifNode, s)
	blockNode := n.right
	if isTrue {
		blockNode = n.left
	}

	execBlockNode(blockNode, s, retVal)
}

func execWhileNode(n *node, s *scope, retVal []string) {
	//eval if in every loop
	for {
		if !evalBoolValue(n.ifNode, s) {
			break
		}

		execBlockNode(n.left, s, retVal)
	}
}

func execFnNode(e *node, s *scope) []*val {
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
func execBlockNode(n *node, parentScope *scope, retTypes []string) []*val {
	blockScope := new(scope)
	blockScope.parent = parentScope
	blockScope.localVars = make(map[string]*val)
	for i := range n.exprs {
		e := n.exprs[i]

		// check if exprs is a return expression
		if e.kind == "fnCall" && e.left.val == "döndür" {
			myPrintln("entered return")
			retval := evalReturnVal(e, blockScope, retTypes)
			return retval
		}

		// forgets children scopes after exec
		def, val := exec(e, blockScope, retTypes)
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
	printFn.GoRef = func(v []*val) []*val {
		myPrintln("WRITE CALLED")
		for i := range v {

			switch v[i].typeName {
			case "tamsayı":
				fmt.Print(v[i].intval, " ")
			case "yapı":
				fmt.Println("yapı:", v[i].objVal.name)
				for k, val := range v[i].objVal.keys {
					fmt.Println(k, *val)
				}
			case "metin":
				fmt.Print("metin:", v[i].strval)
			}
		}
		fmt.Println()
		myPrintln("\nWRITE FINISHED")
		return nil
	}
	printFn.isParamsParametric = true
	printFn.signature = []val{{typeName: ""}}

	printFnVal := new(val)
	printFnVal.typeName = "fn"
	printFnVal.funcVal = printFn
	s.localVars["yazdır"] = printFnVal

}
