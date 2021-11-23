package loopptr

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const category = "loopptr"

type report func(id *ast.Ident)

// Analyzer is an analyzer for checking that addresses of range loop variables are only taken in a safe way.
var Analyzer = &analysis.Analyzer{
	Name: category,
	Doc:  "check that range loop variables do not have their addresses taken inside loop body in a potentially non-final iteration",
	Run:  analyze,
}

func analyze(p *analysis.Pass) (interface{}, error) {
	reporter := func(id *ast.Ident) {
		p.Report(analysis.Diagnostic{
			Pos:      id.Pos(),
			Message:  fmt.Sprintf("taking address of range variable '%v'", id.Name),
			Category: category,
		})
	}
	for _, f := range p.Files {
		ast.Walk(rangeLoopVisitor(reporter), f)
	}
	return nil, nil
}
