package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"mbo/loop-pointer/analysis"
)

func main() {
	singlechecker.Main(analysis.Analyzer)
}
