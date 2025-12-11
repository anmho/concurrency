package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

func main() {

	// create a function that asynchronously results from get requests to a list of urls

	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *Result {
		results := make(chan *Result)

		go func() {
			defer close(results) // you could technically not close it but you need to signal to the consumer that you are done producing

			for _, url := range urls {
				res, err := http.Get(url)
				select {
				case <-done:
					return
				default:
					results <- &Result{Error: err, Response: res}
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
		} else {
			fmt.Printf("Response: %v\n", result.Response.Status)
		}
	}

}
