package rangeloopaddr

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const category = "rangeloopptr"

type reporter func(id *ast.Ident)

// Analyzer is an analyzer for checking that addresses of range loop variables are only taken in a safe way.
// onReport parameter allows a caller to detect if diagnostics have been reported since singlechecker.Main doesn't
// expose that information.
var Analyzer = &analysis.Analyzer{
	Name: category,
	Doc:  "check that addresses of range loop variables aren't taken inside loop body if it may not be the final iteration",
	Run: func(p *analysis.Pass) (interface{}, error) {
		return analyze(p)
	},
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
