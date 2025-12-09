package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})

	go func (done chan interface{}){
		for {
			select {
			case <- done:
				return
			default:
				fmt.Println("do work")
				time.Sleep(3 * time.Second)
			}
		}
	}(done)

	
	// signal completion after 30 seconds

	fmt.Println("cancelling")
	time.Sleep(10 * time.Second)
	fmt.Println("signalling to complete")
	// same thing as
	// done <- struct{}{}
	// the receiver will read a zero value
	close(done)
}