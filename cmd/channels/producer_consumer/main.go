package main

import (
	"fmt"
	"sync"
)

func main() {
	stringStream := make(chan string) // unbuffered channel

	// producer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stringStream <- "hello"
		wg.Done()
	}()

	// consumer. will wait until theres a value in the channel to get
	fmt.Println(<-stringStream)
}
