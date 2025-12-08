package main

import (
	"fmt"
	"sync"
	"time"
)

var readyCond = sync.NewCond(&sync.Mutex{})
var isReady bool = false

func worker(wg *sync.WaitGroup, i int) {
	readyCond.L.Lock()
	defer wg.Done()
	defer readyCond.L.Unlock()

	fmt.Println("worker", i, "waiting on ready")
	for !isReady {
		readyCond.Wait()
	}

	fmt.Println("Worker ", i, " activated!")
}

func setup() {
	readyCond.L.Lock()
	defer readyCond.L.Unlock()

	isReady = true
	readyCond.Broadcast()
}

func main() {
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go worker(&wg, i)
	}

	time.Sleep(5 * time.Second)
	setup()
	fmt.Println("setup done")

	wg.Wait()
}