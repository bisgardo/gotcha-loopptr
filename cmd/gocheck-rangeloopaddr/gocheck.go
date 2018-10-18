package main

import (
	"github.com/halleknast/gocheck-rangeloopaddr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(rangeloopaddr.Analyzer)
}
