package rangeloopaddr

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

const category = "rangeloopptr"

// loopPtrVisitor walks the body of the provided range statement, checking that
// addresses of loop variables are not being taken.
type objPtrVisitor struct {
	pass        *analysis.Pass
	parentObjs  []*ast.Object
	currentObjs []*ast.Object
}

func (v *objPtrVisitor) Visit(node ast.Node) ast.Visitor {
	// If hitting a new range loop, add the loop vars to list of checked objects of the current scope.
	if n, ok := node.(*ast.RangeStmt); ok {
		return &objPtrVisitor{
			pass:        v.pass,
			parentObjs:  v.parentObjs,
			currentObjs: append(loopVarObjs(n), v.currentObjs...),
		}
	}

	// TODO Allow taking address in block that's *guaranteed* to break loop (use ctrlflow pass).

	if _, ok := node.(*ast.ReturnStmt); ok {
		// When recursing into a return statement, exclude objects in current function's scope from check.
		// Objects from parent scopes should still be checked.
		return &objPtrVisitor{
			pass:        v.pass,
			parentObjs:  v.parentObjs,
			currentObjs: nil,
		}
	}
	if _, ok := node.(*ast.FuncLit); ok {
		// Merge parent and current objects to be the parent objects inside nested function.
		var objs []*ast.Object
		objs = append(objs, v.parentObjs...)
		objs = append(objs, v.currentObjs...)
		return &objPtrVisitor{
			pass:        v.pass,
			parentObjs:  objs,
			currentObjs: nil,
		}
	}

	if len(v.parentObjs) == 0 && len(v.currentObjs) == 0 {
		// Skip actual check if there're no objects to test against.
		return v
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

	check(v.pass, id, v.parentObjs, v.currentObjs)
	return v
}

func check(p *analysis.Pass, id *ast.Ident, objses ...[]*ast.Object) {
	// Check for objects in scope of current and parent objects.
	for _, objs := range objses {
		for _, obj := range objs {
			if obj == id.Obj {
				p.Report(analysis.Diagnostic{
					Pos:      id.Pos(),
					Message:  fmt.Sprintf("taking address of range variable '%v'", id.Name),
					Category: category,
				})
			}
		}
	}
}

type rangeLoopVisitor struct {
	pass *analysis.Pass
}

func (v *rangeLoopVisitor) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.RangeStmt)
	if !ok {
		// Node is not a range loop.
		return v
	}

	return &objPtrVisitor{
		pass:        v.pass,
		parentObjs:  nil,
		currentObjs: loopVarObjs(node),
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
