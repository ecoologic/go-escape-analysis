package main

import (
	"os"
	"runtime"
	"runtime/pprof"
)

func getBadData() *[]byte {
	return &[]byte{1, 2, 3, 4, 5}
}

func getGoodData(data *[]byte) { // does not escape
	*data = []byte{1, 2, 3, 4, 5}
}

func main() {
	memFile, err := os.Create("mem.pprof")
	if err != nil {
		panic(err)
	}
	defer memFile.Close()

	run()

	if err = pprof.WriteHeapProfile(memFile); err != nil {
		panic(err)
	}
}

func run() {
	badData := getBadData() // does not escape
	runtime.GC()
	println("getBadData:", *badData, len(*badData))

	var goodData []byte
	getGoodData(&goodData)
	runtime.GC()
	println("getGoodData", goodData, len(goodData))
}

// $ go run main.go
// getBadData: [5/5]0xc00007aec3 5
// getGoodData [5/5]0xc0000200f0 5

// $ go build -gcflags="-m" main.go
// # command-line-arguments
// ./main.go:9:6: can inline getBadData
// ./main.go:13:6: can inline getGoodData
// ./main.go:32:23: inlining call to getBadData
// ./main.go:37:13: inlining call to getGoodData
// ./main.go:18:27: inlining call to os.Create
// ./main.go:26:33: inlining call to pprof.WriteHeapProfile
// ./main.go:26:33: inlining call to pprof.writeHeap
// ./main.go:10:9: &[]byte{...} escapes to heap
// ./main.go:10:16: []byte{...} escapes to heap
// ./main.go:13:18: data does not escape
// ./main.go:14:16: []byte{...} escapes to heap
// ./main.go:32:23: &[]byte{...} does not escape
// ./main.go:32:23: []byte{...} does not escape
// ./main.go:37:13: []byte{...} escapes to heap
