package main

import (
	"fmt"
	"time"
)

// DoWork returns two channels
// 1) a channel for heartbeats
// 2) A generator that streams values from input nums
func DoWork(
	done <-chan interface{},
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {
		// close when we get a done signal
		defer close(heartbeat)
		defer close(intStream)

		// simulate a delay, cpu, network etc
		time.Sleep(2 * time.Second)

		for _, n := range nums {
			select {
			// send heartbeat if theres a
			case heartbeat <- struct{}{}:
			default:
			}

			// send number if
			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeat, intStream
}

// heartbeats that occur at the beginning of a unit of work
func main() {

	done := make(chan interface{})
	heartbeat, intStream := DoWork(done, 1, 2, 3, 4, 5)

	defer close(done)

	for {
		select {
		// have to check that the value we read is not coming from a closed channel
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("heartbeat")
			} else {
				return
			}
		case v, ok := <-intStream:
			if ok {
				fmt.Printf("value %v\n", v)
			} else {
				return
			}
		}
	}
}
