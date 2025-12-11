package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Result struct {
	Response *http.Response
	Error    error
}

func streamGet(
	done <-chan interface{},
	urls []string) <-chan *Result {
	resChan := make(chan *Result)

	go func() {
		defer close(resChan)
		var wg sync.WaitGroup
		wg.Add(len(urls))
		for _, url := range urls {
			go func() {
				res, err := http.Get(url)

				resChan <- &Result{res, err}
				wg.Done()
			}()
		}
		wg.Wait()
	}()

	return resChan
}

func main() {
	done := make(chan interface{})
	stream := streamGet(
		done,
		[]string{
			"https://jsonplaceholder.typicode.com/posts",
			"https://jsonplaceholder.typicode.com/users",
		})
	for result := range stream {
		fmt.Printf("%+v\n", result)
	}
}
