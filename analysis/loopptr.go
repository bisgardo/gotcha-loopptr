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
	keyObj          *ast.Object
	valObj          *ast.Object
	checkReturnStmt bool
}

func (v *objPtrVisitor) Visit(node ast.Node) ast.Visitor {
	if !v.checkReturnStmt {
		if _, ok := node.(*ast.ReturnStmt); ok {
			// Don't descent into return statement if configured to not do that.
			return nil
		}
		if _, ok := node.(*ast.FuncLit); ok {
			// Do check return statements inside function literals.
			return &objPtrVisitor{
				pass:            v.pass,
				keyObj:          v.keyObj,
				valObj:          v.valObj,
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
	obj := id.Obj
	if obj == v.keyObj || obj == v.valObj {
		// Taking address of object corresponding to key or value of range loop.
		v.pass.Report(analysis.Diagnostic{
			Pos:      id.Pos(),
			Message:  fmt.Sprintf("taking address of range variable '%v'", id.Name),
			Category: category,
		})
	}

	return v
}

type rangeLoopVisitor struct {
	pass *analysis.Pass
}

func (v *rangeLoopVisitor) Visit(n ast.Node) (w ast.Visitor) {
	node, ok := n.(*ast.RangeStmt)
	if !ok {
		// Node is not a range loop; descend.
		return v
	}

	key, _ := node.Key.(*ast.Ident)
	val, _ := node.Value.(*ast.Ident)

	if key == nil && val == nil {
		return
	}

	var keyObj *ast.Object
	var valObj *ast.Object
	if key != nil && key.Name != "_" {
		keyObj = key.Obj
	}
	if val != nil && val.Name != "_" {
		valObj = val.Obj
	}

	return &objPtrVisitor{
		pass:            v.pass,
		keyObj:          keyObj,
		valObj:          valObj,
		checkReturnStmt: false,
	}
}
