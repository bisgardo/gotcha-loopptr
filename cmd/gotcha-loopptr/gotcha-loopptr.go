package main

import (
	"github.com/halleknast/gotcha-loopptr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loopptr.Analyzer)
}
