package main

import (
	"github.com/bisgardo/gotcha-loopptr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loopptr.Analyzer)
}
