# gocheck-loopptr

Static checker for Go that reports code where the [address operator](https://golang.org/ref/spec#Address_operators) (`&`)
is used on a loop variable.

## Problem

Taking the address of a loop variable is a fairly common source of bugs.
The reason is that the loop variable(s) are [reused in each iteration](https://golang.org/ref/spec#For_range)
(there's an [open proposal](https://github.com/golang/go/issues/20733) for changing these semantics in "Go 2").
This is the case for both "regular" for loops and range loops, but is probably less surprising for the regular loops.
This tool currently only checks range loops, but inclusion of regular loops is being considered.

The rationale for assuming that such code is buggy is that the programmer probably assumes
that the resulting pointer is different in each iteration - just like the contents of the variables are.
If not, there are no cases where it couldn't be trivially rewritten to a "safe" style that's at least as efficient
(see "analysis" section below).

[Example](https://play.golang.org/p/64R8BC_egcz):

```
var ps []*int
vs := []int{1, 2, 3}
for _, v := range vs {
	ps = append(ps, &v)
}

// Might expect ps[i] to point to vs[i],
// but they all point to the last element of vs (3).
for i, p := range ps {
	println(i, ":", *p)
}
```

## Analysis

The project is built using the [analysis API](https://godoc.org/golang.org/x/tools/go/analysis)
and includes a standalone tool for invoking it.

The tool reports code such as the above like so:

```
for _, v := range vs {
	ps = append(ps, &v)
}
```

[TODO: Add reporting output and examples - including return statements and examples of rewriting false positives]

## Install tool

```
go get github.com/halleknast/gocheck-loopptr/cmd/gocheck-loopptr
cd <project root>
gocheck-loopptr ./...
```

No output means "nothing to report".
