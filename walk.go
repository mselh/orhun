package main

type program struct {
	rootNode  *node // root ast node
	rootScope *scope
}

type scope struct {
	parent    *scope
	localVars map[string]any
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

	// start from entry
	entryNode := p.rootNode.left
	entryScope := new(scope)
	entryScope.parent = p.rootScope
	entryScope.localVars = make(map[string]any)
	for i := range entryNode.exprs {
		e := entryNode.exprs[i]
		// forgets children scopes after exec
		def, val := exec(e, entryScope)
		if def != nil {
			entryScope.localVars[def] = val
		}

	}

}

// execs a node
// if it s a new node, returns that new node
// else only executes it.
func exec(e *node, parentScope *scope) {
	if e.kind == "" {
		// .....
	}
}

// helpers
func newProgram(root *node) *program {
	rootScope := new(scope)
	rootScope.localVars = make(map[string]any)

	return &program{
		rootNode:  root,
		rootScope: rootScope,
	}
}
