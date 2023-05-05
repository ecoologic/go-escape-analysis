package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
)

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

func bigData() []byte {
	file, err := os.Open("./2023-values.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the file into memory
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return data
}

func run() {
	badData := getBadData()
	runtime.GC()
	println("getBadData:", *badData, len(*badData))

	var goodData []byte
	getGoodData(&goodData)
	runtime.GC()
	println("getGoodData", goodData, len(goodData))
}

// More complex methods wouldn't get inlined
// I'm keeping it simple and disabling the compiler optimisation
//go:noinline
func getBadData() *[]byte {
	data := bigData() // :52:2: data escapes to heap
	return &data
}

//go:noinline
func getGoodData(data *[]byte) { // :57:18: data does not escape
	*data = bigData()
}
