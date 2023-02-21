package main

import (
	"fmt"
	"net/http"
)

func ErrorHandleMain() {
	done := make(chan any)
	defer close(done)

	urls := []string{"https://www.google.com", "https://www.youtube.com", "https://bad"}
	eh := ErrorHandle{}
	for res := range eh.CheckStatus(done, urls) {
		if res.Error != nil {
			fmt.Println(res.Error)
			continue
		}

		fmt.Println(res.Response.Status)
	}
}

type ErrorHandle struct{}

func (e ErrorHandle) CheckStatus(done <-chan any, urls []string) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)
		for _, url := range urls {
			res, err := http.Get(url)
			result := Result{Response: res, Error: err}

			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()

	return results
}

type Result struct {
	Response *http.Response
	Error    error
}
