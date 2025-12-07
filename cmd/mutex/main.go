package main

import (
	"fmt"
	"sync"
)

func increment(count *int, m *sync.Mutex) {
	m.Lock()
	defer m.Unlock()
	*count++
}

func decrement(count *int, m *sync.Mutex) {
	m.Lock()
	defer m.Unlock()
	*count--
}

func main() {
	for range 50 {
		var arithmetic sync.WaitGroup

		var count int

		var m sync.Mutex
		for i := 0; i <= 5; i++ {
			arithmetic.Add(1)

			go func() {
				defer arithmetic.Done()
				increment(&count, &m)
			}()
		}

		for i := 0; i <= 5; i++ {
			arithmetic.Add(1)
			go func() {
				defer arithmetic.Done()
				decrement(&count, &m)
			}()
		}

		arithmetic.Wait()

		fmt.Println("count", count)

	}

}
