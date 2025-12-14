package main

import (
	"fmt"
	"math/rand"
)

// heartbeats that occur at the beginning of a unit of work
func main() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		// one buffer so we can store a ping. why?
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)

		go func() {

			// close channels when we are notified to exit
			defer close(heartbeatStream)
			defer close(workStream)

			for i := 0; i < 10; i++ {

				// heartbeat on each piece of work
				// if theres already there but not consumed, skip
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}

				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	heartbeat, results := doWork(done)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}
