package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	// say hello and goodbye concurrently

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()
}

func printGreeting(done chan interface{}) error {
	greeting, err := genGreeting(done)

	if err != nil {
		return err
	}

	fmt.Printf("%v world!\n", greeting)
	return nil

}

func printFarewell(done chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	fmt.Printf("%v world!\n", farewell)

	return nil
}

func genGreeting(done chan interface{}) (string, error) {
	// locale, err := locale(done)
	// if err != nil {
	// 	return "", err
	// }
	// switch locale {
	// case "EN/US":
	// 	return "Hello", nil
	// default:
	// 	return "", errors.New("unsupported locale")
	// }

	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "Hello", nil
	default:
		return "", errors.New("unsupported locale")
	}
}

func genFarewell(done chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "Goodbye", nil
	default:
		return "", errors.New("unsupported locale")
	}
}

func locale(done chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(1 * time.Minute): // auto cancel after 1 minute
	}
	return "EN/US", nil
}
