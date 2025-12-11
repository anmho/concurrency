package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	var clickRegistered sync.WaitGroup

	subscribe := func(c *sync.Cond, fn func()) {
		// we wait until the goroutine is "running" because that means it successfully registerd
		// once its running then it waits for the cond event broadcast
		var goroutineRunning sync.WaitGroup

		func() {
			goroutineRunning.Add(1)
			c.L.Lock()
			defer c.L.Unlock()

			c.Wait() // wait for event to fire
			// call the callback
			fn()
		}()

		goroutineRunning.Wait()
	}

	buttonClicked := sync.NewCond(&sync.Mutex{})
	subscribe(buttonClicked, func() {
		fmt.Println("Maximizing window")
		clickRegistered.Done()
	})

	subscribe(buttonClicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})

	subscribe(buttonClicked, func() {
		fmt.Println("Display update version toast")
		clickRegistered.Done()
	})

	// wait for all click event handlers to be registered

	// at this point they are all registered and waiting

	// simulate button click
	fmt.Println("simulating button click event")
	button.Clicked.Broadcast()

	clickRegistered.Wait()

	fmt.Println("button click")
}
