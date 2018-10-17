package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"mbo/gocheck-rangeloopaddr/rangeloopaddr"
)

func main() {
	singlechecker.Main(rangeloopaddr.Analyzer)
}
