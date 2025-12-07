package main

import (
	"sync"
	"time"
)

type value struct {
	mu    sync.Mutex
	value int
}

func main() {
	var wg sync.WaitGroup

	printSum := func(v1, v2 *value) {
		defer wg.Done()

		// acquire lock1
		v1.mu.Lock()

		// release lock 1 (push onto the cleanup stack, lifo)
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)

		v2.mu.Lock()
		// release lock 2(push onto the lceanup stack, lifo)
		defer v2.mu.Unlock()
	}

	var a, b value

	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Add(2)
	wg.Wait()
}
