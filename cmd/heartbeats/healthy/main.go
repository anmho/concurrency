package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(
		done chan interface{},
	) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		
		// sends heartbeat to heartbeat channel while not done
		go func() {
			// close channels when we are told to stop
			// we will not stop otherwise
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(1 * time.Second)
			workGen := time.Tick(2 * time.Second)


			// generate one heartbeat
			// we can skip generating if noone is listening 
			// logically the heartbeat doesn't matter if there is noone there to listen
			// “If a tree falls in a forest and no one is around to hear it, does it make a sound?”
			sendPulse := func() {
				select {
				case heartbeat <-struct{}{}:
				default: // if heartbeat was read/consumed, then send one else do skip
				}
			}
			

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done: // we need this if sub routine must be pre-emptable
						return
					case results <- r:
						return
					case <- pulse: // we need this so w can send pulse even if there is no listener/result chan is full
						sendPulse() // we were not returning when sending a result
					}
				}
			}


			// check which channels are ready for action
			for {
				select {
				case <-done: // check for done
					return
				case <- pulse: // check for pulse
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()

		return heartbeat, results
	}

	done := make(chan interface{})
	heartbeat, results := doWork(done)
	time.AfterFunc(20 * time.Second, func() { close(done) }) 
	for {
		select {
		case _, ok := <- heartbeat:
			if !ok { // return if channel was closed (e.g. done signal was sent)
				return
			} else {
				fmt.Println("heartbeat")
			}
			
		case r, ok := <- results:
			if !ok { // value was returned from a closed channel which means we are done
				return
			} else {
				fmt.Println("result", r)
			}
			fmt.Println("result", r)
		case <- time.After(10 * time.Second):
			fmt.Println("timed out with no result or heartbeat")
			return
		}
	}
}