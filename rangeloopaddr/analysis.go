package rangeloopaddr

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "lock",
	Doc:  "check that addresses of range loop variables aren't taken inside loop body",
	Run:  analyze,
}

func analyze(p *analysis.Pass) (interface{}, error) {
	for _, f := range p.Files {
		ast.Walk(&rangeLoopVisitor{pass: p}, f)
	}
	return nil, nil
}
