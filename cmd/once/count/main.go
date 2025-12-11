package main

import (
	"fmt"
	"sync"
)

var count int

func increment() {
	count++
}

func main() {
	var once sync.Once

	// increment 100 times concurrently
	var wg sync.WaitGroup

	routines := 100
	wg.Add(routines)
	for range routines {
		go func() {
			// increment
			once.Do(increment)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("final count", count)
}
