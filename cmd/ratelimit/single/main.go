package main

import (
	"context"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

const (
	PerSecond = 1
	MaxBurst  = 5
)

func Open() *APIConnection {
	return &APIConnection{
		// 1 per second, burst of 1
		rateLimiter: rate.NewLimiter(rate.Limit(PerSecond), MaxBurst),
	}
}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// blocking rate limit instead of fast fail
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}

	return nil
}
func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// blocking rate limit instead of fast fail
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

// heartbeats that occur at the beginning of a unit of work
func main() {
	// set log flags
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	// open api connection
	apiConnection := Open()

	var wg sync.WaitGroup

	// add 20 to the task counter
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}
