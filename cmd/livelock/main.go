package main

import (
	"bytes"
	"sync"
	"time"
)

func main() {

	cadence := sync.NewCond(&sync.Mutex{})

	// run in the background persistently until program shutdown
	go func() {

		// every 1 millisecond, signal that its okay to go for any goroutines waiting on cond

		// convenience func to return a read only channel. we are not using the value
		for range time.Tick(1*time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}






}


