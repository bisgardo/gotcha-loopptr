package analysis

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
	pass            *analysis.Pass
	objs            []*ast.Object
	checkReturnStmt bool
}

func (v *objPtrVisitor) Visit(node ast.Node) ast.Visitor {
	// TODO If hitting a new range loop, add the loop vars to list of checked objects.

	if !v.checkReturnStmt {
		if _, ok := node.(*ast.ReturnStmt); ok {
			// Don't descend into return statement if it's in the same function as the loop.
			return nil
		}
		if _, ok := node.(*ast.FuncLit); ok {
			// Do check return statements inside function literals.
			// TODO Might as well remove objects of shadowed variables.
			return &objPtrVisitor{
				pass:            v.pass,
				objs:            v.objs,
				checkReturnStmt: true,
			}
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

	for _, obj := range v.objs {
		if obj == id.Obj {
			// Taking address of object corresponding to key or value variable of range loop.
			v.pass.Report(analysis.Diagnostic{
				Pos:      id.Pos(),
				Message:  fmt.Sprintf("taking address of range variable '%v'", id.Name),
				Category: category,
			})
		}
	}
	return v
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

	objs := loopVarObjs(node)
	if len(objs) == 0 {
		return nil
	}
	return &objPtrVisitor{
		pass:            v.pass,
		objs:            objs,
		checkReturnStmt: false,
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
