package main

import (
	"fmt"
	"net/http"
)

var urls = []string{
	"https://stackoverflow.com/",
	"https://github.com/",
	"https://www.linkedin.com/",
	"http://medium.com/",
	"https://golang.org/",
	"https://www.udemy.com/",
	"https://www.coursera.org/",
	"https://wesionary.team/",
}

func main() {
	c := make(chan string)

	for _, url := range urls {
		go sendRequestAsync(url, c)
	}

	for i := range urls {
		msg := <-c
		fmt.Printf("message %d: %s \n", i, msg)
	}
}

func sendRequestAsync(url string, c chan string) {
	_, err := http.Get(url)
	if err != nil {
		c <- url + "responded with error"
	}
	c <- url + "is up"
}
