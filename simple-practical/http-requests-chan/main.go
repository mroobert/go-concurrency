package main

// NOTES
// 1.Once a channel is closed, you can't send further values on it else it panics.

// 2.You can close the channel only once (attempting to close an already closed channel also panics).
// And you should do it when all the goroutines that send values on it are done.
// In order to do this, you need to detect when all the sender goroutines are done. An idiomatic way to detect this is to use sync.WaitGroup.

// 3.The WaitGroup.Wait() will wait until all sender goroutines are done, and only after this and only once will it close the channel.
// We want to detect this "global" done event and close the channel while processing the values sent on it is in progress, so we have to do this in its own goroutine.

// 4.The for ... range will run until the channel is closed.
// And since it runs in the main goroutine, the program will not exit until all the values are properly received and processed from the channel.
// The for ... range construct loops until all the values are received(the ones were sent before the channel was closed).

// !!!!!
//! Note that it's important that the for ... range runs in the main goroutine, and the code which waits for the WaitGroup and closes the channel is in its own (non-main) goroutine; and not the other way around.
//! If you would switch these 2, it might cause an "early exit", that is not all values might be received and processed from the channel.
//! The reason for this is because a Go program exits when the main goroutine finishes (spec: Program execution).
//! It does not wait for other (non-main) goroutines to finish.
//! So if waiting and closing the channel would be in the main goroutine, after closing the channel the program could exit at any moment,
//! not waiting for the other goroutine that in this case would loop to receive values from the channel.

import (
	"fmt"
	"net/http"
	"sync"
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

var wg sync.WaitGroup

func main() {
	c := make(chan string)

	fmt.Println("Starting sending requests...")
	for _, url := range urls {
		wg.Add(1)
		go sendRequestAsync(url, c)
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	for message := range c {
		fmt.Printf("message: %s \n", message)
	}

	fmt.Println("End...")
}

func sendRequestAsync(url string, c chan<- string) {
	defer wg.Done()
	_, err := http.Get(url)
	if err != nil {
		c <- url + "responded with error"
	}
	c <- url + "is up"
}
