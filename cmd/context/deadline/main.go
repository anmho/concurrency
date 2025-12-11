package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	// done := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // trigger the cancel when we are done. indicate to goroutines that we should clean up or are done.
	// defer close(done)

	// say hello and goodbye concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)

	if err != nil {
		return err
	}

	fmt.Printf("%v world!\n", greeting)
	return nil

}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%v world!\n", farewell)

	return nil
}

func genGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel() // this will send to done chan after 1 second. defer cancel actually CANCELS the timeout after we cleanup
	// with timeout will send to done chan after 1 second

	
	switch locale, err := locale(ctx); {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "Hello", nil
		default:
			return "", errors.New("unsupported locale")
	}
}

func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "Goodbye", nil
		default:
			return "", errors.New("unsupported locale")
	}
}


func locale(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok { // fail fast
		// if the deadline is before the completion time, then we won't finish
		completionTime := time.Now().Add(1 * time.Minute)
		if completionTime.After(deadline) {
			return "", context.DeadlineExceeded
		}
	}


	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1*time.Minute): // return after 1 minute
	}
	return "EN/US", nil
}
