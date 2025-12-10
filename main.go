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

	

	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)

	for i := range multiply(done, intStream, 3) {
		fmt.Println(i)
	}
}