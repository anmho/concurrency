package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func backgroundJob(ctx context.Context) {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			fmt.Println("tick")
		case <-ctx.Done():
			fmt.Println("boom")
			return
		}
	}
}


func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	
	go func () {
		defer wg.Done()
		backgroundJob(ctx)
	}()
	wg.Wait()
}