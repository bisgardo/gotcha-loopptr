package rangeloopaddr

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer checks that addresses of range loop variables are only taken in a safe way.
var Analyzer = &analysis.Analyzer{
	Name: category,
	Doc:  "check that addresses of range loop variables aren't taken inside loop body if it may not be the final iteration",
	Run:  analyze,
}

func analyze(p *analysis.Pass) (interface{}, error) {
	for _, f := range p.Files {
		ast.Walk(&rangeLoopVisitor{pass: p}, f)
	}
	return nil, nil
}
