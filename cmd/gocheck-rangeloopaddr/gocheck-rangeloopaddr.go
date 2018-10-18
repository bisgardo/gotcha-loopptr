package main

import (
	"os"

	"github.com/halleknast/gocheck-rangeloopaddr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	fail := false
	singlechecker.Main(rangeloopaddr.NewAnalyzer(func() {
		fail = true
	}))
	if fail {
		os.Exit(2)
	}
}
