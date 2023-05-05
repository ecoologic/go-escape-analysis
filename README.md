# Go Escape Analysis

When manipulating large or variable-sized data in Go functions, it is important to understand the performance implications of passing data as a value versus passing a pointer to the data. By profiling a small program and using escape analysis, we can demonstrate why passing a pointer is often more efficient.

If you declare a local pointer inside a function and don't return it or store it somewhere else that is accessible outside of the function, the garbage collector can free the related memory when the function is over. But if you return a pointer, the memory will be freed by the garbage collector only when there are no more references to it in the program.

[From Golang FAQ](https://tip.golang.org/doc/faq#stack_or_heap)

> if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors.


Prefer:

```go
func getData(data *[]byte) {
	*data = bigData()
}
```

over:

```go
func getData() *[]byte {
	data := bigData()
	return &data
}
```

This is why we need to use readers like so:

```go
buffer := make([]byte, 1024)
n, err := file.Read(buffer)
```

instead of just returning the buffer from `file.Read`.

* [Understanding Allocations: the Stack and the Heap](https://youtu.be/ZMZpH4yT7M0)
* [Escape Analysis and Memory Profiling](https://youtu.be/2557w0qsDV0)

## Commands

```sh
go run main.go -memprofile mem.out
go build -gcflags="-m -m" main.go
go tool pprof -alloc_space mem.pprof
```

## Run it!

```sh
~/dev/go/escape-analysis(main)✗$ go version
go version go1.19.5 darwin/amd64

~/dev/go/escape-analysis(main)✗$ go run main.go

getBadData: [593966/704512]0xc00038e000 593966
getGoodData [593966/704512]0xc000580000 593966

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

## PProf

```sh
$ go tool pprof -alloc_space mem.pprof
Type: alloc_space
Time: May 5, 2023 at 12:55pm (AEST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list getBadData
Total: 13.13MB
ROUTINE ======================== main.getBadData in ~/dev/go/escape-analysis/main.go
         0     6.60MB (flat, cum) 50.30% of Total
         .          .     50://go:noinline
         .     6.60MB     51:func getBadData() *[]byte {
         .          .     52:	data := bigData()
         .          .     53:	return &data
         .          .     54:}
         .          .     55:
         .          .     56://go:noinline
```
