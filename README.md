# Go Escape Analysis

Profiling a small program and using escape analysis to demonstrate why, when manipulating big or variable size data, we should pass a pointer to it in Go functions:

Prefer:

```go
func getData(data *[]byte) {
	*data = bigData()
}
```

over this:

```go
func getData() *[]byte {
	data := bigData()
	return &data
}
```

[Dig deeper here](https://youtu.be/2557w0qsDV0).

Commands used:

```sh
go run main.go
go build -gcflags="-m -m" main.go
go tool pprof -alloc_space mem.pprof
```

How it runs:

```sh
~/dev/go/escape-analysis(main)✗$ go version
go version go1.19.5 darwin/amd64

~/dev/go/escape-analysis(main)✗$ go run main.go

getBadData: [1317914/1400832]0xc000220000 1317914
getGoodData [1317914/1400832]0xc000700000 1317914

~/dev/go/escape-analysis(main)✗$ go build -gcflags="-m -m" main.go
# command-line-arguments
./main.go:24:6: cannot inline bigData: unhandled op DEFER
./main.go:25:22: inlining call to os.Open
./main.go:32:29: inlining call to ioutil.ReadAll
./main.go:51:6: cannot inline getBadData: marked go:noinline
./main.go:57:6: cannot inline getGoodData: marked go:noinline
./main.go:39:6: cannot inline run: function too complex: cost 258 exceeds budget 80
./main.go:10:6: cannot inline main: unhandled op DEFER
./main.go:11:27: inlining call to os.Create
./main.go:19:33: inlining call to pprof.WriteHeapProfile
./main.go:19:33: inlining call to pprof.writeHeap
./main.go:29:18: bigData capturing by value: .autotmp_14 (addr=false assign=false width=8)
./main.go:52:2: data escapes to heap:
./main.go:52:2:   flow: ~r0 = &data:
./main.go:52:2:     from &data (address-of) at ./main.go:53:9
./main.go:52:2:     from return &data (return) at ./main.go:53:2
./main.go:52:2: moved to heap: data
./main.go:57:18: data does not escape
./main.go:15:21: main capturing by value: .autotmp_12 (addr=false assign=false width=8)
```

`getBadData` escapes to the heap:

```
./main.go:52:2: data escapes to heap:
```

`getGoodData` doesn't:

```
./main.go:57:18: data does not escape
```

Note on `//go:noinline`: This disables the compiler optimisation, for more complex methods we would read something like this:

```
./main.go:10:6: cannot inline main: unhandled op DEFER
```

Note on `mem.pprof`: This is for deeper analysis.
