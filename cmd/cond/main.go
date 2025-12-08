package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// a way to do event driven waiting

func exampleBroadcast() {
	// with sync.Cond we can c.Signal() to wake up one of the goroutines. This is nondeterministic
	// or we can c.Broadcast() to wake up goroutines waiting on the cond

	
	readyCond := sync.NewCond(&sync.Mutex{})


	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			readyCond.L.Lock()
			defer readyCond.L.Unlock()
			fmt.Println("worker", i, "waiting for condition")
			readyCond.Wait() // this w

			waitTime := time.Duration(1 + 5 * rand.Float64()) * time.Second

			fmt.Println("worker", i, "doing work for", waitTime)
			time.Sleep(waitTime)
			fmt.Println("worker", i, "done")

			fmt.Println()
			wg.Done()
		}()
	}

	
	fmt.Println("waiting for condition to be met")
	time.Sleep(1 *time.Second)
	fmt.Println("condition was met, notifying and beginning work")

	// some condition is met
	readyCond.Broadcast()

	wg.Wait()
	fmt.Println("Work completed")
}

func exampleSignal() {
	// with sync.Cond we can c.Signal() to wake up one of the goroutines. This is nondeterministic
	// or we can c.Broadcast() to wake up goroutines waiting on the cond

	readyCond := sync.NewCond(&sync.Mutex{})


	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			readyCond.L.Lock()
			defer readyCond.L.Unlock()
			fmt.Println("worker", i, "waiting for condition")
			readyCond.Wait() // this w

			waitTime := time.Duration(1 + 5 * rand.Float64()) * time.Second

			fmt.Println("worker", i, "doing work for", waitTime)
			time.Sleep(waitTime)

			fmt.Println()
			wg.Done()
		}()
	}



	
	fmt.Println("waiting for condition to be met")
	time.Sleep(10 *time.Second)
	fmt.Println("condition was met, notifying and beginning work")

	// some condition is met
	// lets do them one at a time

	for range 5 {
		readyCond.Signal()
	}

	wg.Wait()
}

func main() {
	fmt.Println("Running broadcast example")
	exampleBroadcast()
	// fmt.Println("Running signal example")
	// exampleSignal()
}
