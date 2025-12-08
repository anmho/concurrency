package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)


	// createa buffered chan of 4
	intStream := make(chan int, 4)

	// fire off gorotutine
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	// gets items until channel is closed. then breaks the loop
	for i := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v", i)
	}
}package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)


	// createa buffered chan of 4
	intStream := make(chan int, 4)

	// fire off gorotutine
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	// gets items until channel is closed. then breaks the loop
	for i := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v", i)
	}
}