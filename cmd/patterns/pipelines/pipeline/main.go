package main

import "fmt"

func main() {

	generator := func(done <- chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func () {
			defer close(intStream) // notify the consumer that we are done producing
			for i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func (done <- chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			// close once we are done producing
			defer close(multipliedStream)

			for i := range intStream {
				select {
				case <-done: // stop if cancelled early
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	add := func (done <-chan interface{}, intStream <- chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream) // signal done producing

			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()

		return addedStream
	}


	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)

	pipeline := add(
		done,
		multiply(done, intStream, 3),
		3,
	)

	for i := range pipeline {
		fmt.Println(i)
	}
}
