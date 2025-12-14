package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	PerSecond = 1
	MaxBurst  = 5
)
type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	// 1/minute
	// the duration it takes to refresh a single token
	return rate.Every(duration / time.Duration(eventCount))
}

// Lets us combine rate limits 
type MultiLimiter struct {
	limiters []RateLimiter
}

func (l *MultiLimiter) Wait(ctx context.Context) error {
	// you wait until in order of most restrictive rate limit
	// this lets you define rate limits by second, my minute, by day etc
	for _, rateLimiter := range l.limiters {
		if err := rateLimiter.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *MultiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func NewMultiLimiter(limiters ...RateLimiter) *MultiLimiter {
	// sort rate limiters by limit (ops/second)
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &MultiLimiter{limiters: limiters}
}

func Open() *APIConnection {
	// 2 per second. 1 op burst
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
	// 10 per minute, 10 burst
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)

	var _ (RateLimiter) = secondLimit
	return &APIConnection{
		// 1 per second, burst of 1
		rateLimiter: NewMultiLimiter(secondLimit, minuteLimit),
	}
}

type APIConnection struct {
	rateLimiter RateLimiter
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
