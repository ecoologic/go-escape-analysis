package main

import (
	"fmt"
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
	file, err := os.Open("./agile.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return data
}

func logGC(data *[]byte, name string) {
	runtime.SetFinalizer(data, func(interface{}) {
		fmt.Printf("%s garbage collected\n", name)
	})
}

func run() {
	badData := getBadData()
	fmt.Println("1\n")
	runtime.GC()
	fmt.Println("2\n")
	println("getBadData:", *badData, len(*badData))

	var goodData []byte
	fmt.Println("3\n")
	getGoodData(&goodData)
	fmt.Println("4\n")
	runtime.GC()
	println("getGoodData", goodData, len(goodData))
	runtime.GC()
	fmt.Println("5\n")
}

//go:noinline
func getBadData() *[]byte {
	data := bigData() // :52:2: data escapes to heap
	logGC(&data, "badData")
	return &data
}

//go:noinline
func getGoodData(data *[]byte) { // :57:18: data does not escape
	*data = bigData()
	logGC(data, "goodData")
}
