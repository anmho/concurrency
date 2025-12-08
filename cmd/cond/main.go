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



func producerConsumer() {
	var capacity = 2
	// initialize as empty with initial capacity of two
	var queue = make([]interface{}, 0, capacity)
	var c = sync.NewCond(&sync.Mutex{})

	popQueue := func() {
		time.Sleep(500 * time.Millisecond)

		c.L.Lock()
		defer c.L.Unlock()


		// if there is space in the queue
		if len(queue) > 0 {
			// remove first item in queue
			queue = queue[1:]

			// do something with it
			time.Sleep(1 *time.Second)
			fmt.Println("Removed item. Queue length:", len(queue))
		}

		c.Signal()
	}


	for range 5 {
		c.L.Lock()

		// we have to guard against spurious wakeups
		for len(queue) == capacity { // full
			fmt.Println("Producer waiting: queue full.")
			c.Wait()
		}


		fmt.Println("Producer adding item.")

		// produce item
		queue = append(queue, struct{}{})

		// add a worker to do work on the new item
		go popQueue()
		
		c.L.Unlock()
	}



}

func broadcastExample() {

}

func main() {

	
	
}
