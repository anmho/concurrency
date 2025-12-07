package main

import (
	"fmt"
	"runtime"
	"sync"
)


func main() {

	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	// read only channel
	var c <- chan interface{}


	var wg sync.WaitGroup

	// mark done but block forever
	noop := func() {
		wg.Done() 
		<- c
	}

	const numGoroutines uint64 = 1e4

	before := memConsumed()


	wg.Add(int(numGoroutines))

	for range numGoroutines {
		go noop()
	}

	wg.Wait()

	after := memConsumed()

	memUsg := after-before
	// in bytes
	fmt.Println("memUsg", memUsg)
	fmt.Printf("%.3fkb", float64(memUsg)/float64(numGoroutines)/1000)

}