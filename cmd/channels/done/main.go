package main

import (
	"fmt"
	"time"
)

func main() {
	// launches a goroutine that does work and returns a channel that tells us when its cancelled
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})

		go func() {
			fmt.Println("doing work")
			defer fmt.Println("doWork exited.")
			defer close(terminated)

			for {
				// none are higher priority than the other
				// blocks if none of the channels have anything
				select {
				case s := <-strings:
					fmt.Println("doing string")
					fmt.Println(s)
				case <-done: // complete when given the signal
					fmt.Println("done")
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan interface{})

	stringsChan := make(chan string, 4)
	stringsChan <- "a"
	stringsChan <- "b"
	stringsChan <- "c"

	terminated := doWork(done, stringsChan)
	go func() {
		time.Sleep(5 * time.Second)
		close(done) // or done <- struct{}{}
	}()
	<-terminated

}
