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
	defer cancel()
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
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("canceled")
	case <-time.After(1 * time.Minute): // auto cancel after 1 minute
	}
	return "EN/US", nil
}
