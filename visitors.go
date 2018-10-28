package rangeloopaddr

import (
	"go/ast"
	"go/token"
)

// objPtrVisitor is a visitor for walking the body of a range statement and
// reporting if addresses of checked loop variables are being taken.
type objPtrVisitor struct {
	reporter    reporter
	parentObjs  []*ast.Object // AST objects to check from the parent scope.
	currentObjs []*ast.Object // AST objects to check from the current function's scope.
}

func (v *objPtrVisitor) Visit(node ast.Node) ast.Visitor {
	// If hitting a new range loop,
	// add the loop vars to list of checked objects of the current scope.
	if n, ok := node.(*ast.RangeStmt); ok {
		return &objPtrVisitor{
			reporter:    v.reporter,
			parentObjs:  v.parentObjs,
			currentObjs: append(loopVarObjs(n), v.currentObjs...),
		}
	}

	if _, ok := node.(*ast.ReturnStmt); ok {
		// When recursing into a return statement,
		// exclude objects in current function's scope from check.
		// Objects from parent scopes should still be checked.
		// Optimization: If there's no parent objects to check,
		// go back to using the lighter visitor rangeLoopVisitor.
		if len(v.parentObjs) == 0 {
			return rangeLoopVisitor(v.reporter)
		}
		return &objPtrVisitor{
			reporter:    v.reporter,
			parentObjs:  v.parentObjs,
			currentObjs: nil,
		}
	}
	if _, ok := node.(*ast.FuncLit); ok {
		// Merge parent and current objects to be the parent scope
		// objects inside nested function.
		var objs []*ast.Object
		objs = append(objs, v.parentObjs...)
		objs = append(objs, v.currentObjs...)
		return &objPtrVisitor{
			reporter:    v.reporter,
			parentObjs:  objs,
			currentObjs: nil,
		}
	}

	ue, ok := node.(*ast.UnaryExpr)
	if !ok || ue.Op != token.AND {
		// Not an "address of" expression.
		return v
	}

	id, ok := ue.X.(*ast.Ident)
	if !ok {
		// Not taking address of identifier.
		return v
	}

	checkIdent(v.reporter, id, v.parentObjs, v.currentObjs)
	return v
}

func checkIdent(r reporter, id *ast.Ident, objses ...[]*ast.Object) {
	// Check for objects in scope of current and parent objects.
	for _, objs := range objses {
		for _, obj := range objs {
			if obj == id.Obj {
				r(id)
			}
		}
	}
}

// rangeLoopVisitor is a light weight visitor for walking an AST node
// and delegating range statement nodes to objPtrVisitor, which does the actual checking.
type rangeLoopVisitor reporter

func (v rangeLoopVisitor) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.RangeStmt)
	if !ok {
		// Node is not a range loop.
		return v
	}

	objs := loopVarObjs(node)
	if len(objs) == 0 {
		// No objects to check.
		return v
	}
	return &objPtrVisitor{
		reporter:    reporter(v),
		parentObjs:  nil,
		currentObjs: objs,
	}
}

func loopVarObjs(node *ast.RangeStmt) []*ast.Object {
	var objs []*ast.Object
	key, _ := node.Key.(*ast.Ident)
	val, _ := node.Value.(*ast.Ident)
	if key != nil && key.Name != "_" {
		objs = append(objs, key.Obj)
	}
	if val != nil && val.Name != "_" {
		objs = append(objs, val.Obj)
	}
	return objs
}
