# gotcha-loopptr

A simple linting tool build on top of the [analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) package.

It checks code for range loops where the loop variables have their addresses taken inside the body.
Unless this happens in a `return` statement, the analyzer will conclude that the address may escape the loop.
As this is a [common source of bugs](https://github.com/golang/go/issues/20733),
the code is flagged as such.

*Example:*

```
for i := range v {
	f(i)      // ok
	g(&i)     // rejected
	return &i // ok
}
```

See [`testdata`](testdata/src/a/rangeloopaddr.go) for more complex examples.

## Build

```
go build ./cmd/gotcha-loopptr
```

## Run

From the root of the project to analyze:

```
gotcha-loopptr ./...
```
