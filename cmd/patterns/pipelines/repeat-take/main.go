package main

import (
	"fmt"
	"math/rand"
)

func main() {

	// repeat repeats the values passed until you tell it to stop
	// repeat := func(
	// 	done <-chan interface{},
	// 	values ...interface{},
	// ) <-chan interface{} {
	// 	repeated := make(chan interface{})

	// 	go func() {
	// 		defer close(repeated) // signal to consumer that we are done
	// 		for {
	// 			for _, v := range values {
	// 				select {
	// 				case <-done:
	// 					return
	// 				case repeated <- v:
	// 				}
	// 			}
	// 		}
	// 	}()
	// 	return repeated
	// }

	// take first n values of input stream and pass to the output channel
	take := func(
		done <-chan interface{},
		valueStream <- chan interface{},
		num int,
	) <- chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream: // receive value from valuestream and send it to takestream if BOTH are ready
				}
			}
		}()
		return takeStream
	}

	repeatFn := func(
		done <- chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select{
				case <- done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}


	done := make(chan interface{})
	defer close(done)

	// go func() {
	// 	defer close(done)
	// 	defer fmt.Println("done")
	// 	time.Sleep(2 * time.Second)
	// }()

	randFn := func() interface{} {
		return rand.Int()
	}

	for v := range take(done, repeatFn(done,  randFn), 2) {
		fmt.Println(v)
	}
}