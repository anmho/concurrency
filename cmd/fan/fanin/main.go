package main

import (
	"fmt"
	"sync"
	"time"
)

func makeChan(
	done chan interface{},
	values ...int) <-chan interface{} {
	c := make(chan interface{})

	go func() {
		defer close(c)

		for _, v := range values {
			select {
			case <-done:
				return
			default:
				time.Sleep(1 * time.Second)
				c <- v
			}
		}
	}()

	return c
}
func main() {
	// fan in merges all the channels into a single channel
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexed := make(chan interface{})
		// we are done when all the reads complete

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()

			for v := range c {
				select {
				case <-done:
					return
				case multiplexed <- v: // the for loop would just loop forever if not cancelled. this way we stop on closed channel from producer chans
				}
			}
		}

		for _, c := range channels {
			wg.Add(1)
			go multiplex(c)
		}

		go func() {
			defer close(multiplexed)
			wg.Wait()
		}()

		return multiplexed
	}

	done := make(chan interface{})
	defer close(done)
	for v := range fanIn(done,
		makeChan(done, 1, 2, 3, 4, 5),
		makeChan(done, 10, 20, 30, 40, 50),
	) {
		fmt.Println(v)
	}
}
