package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

// time go run main.go github.com google.com wikipedia.com youtube.com medium.com bandcamp.com bstn.com coursera.org twitter.com facebook.com soundcloud.com gmail.com udemy.com

// 1. sync = 7s
// 2. async = 2s

var wg sync.WaitGroup
var mut sync.Mutex

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url1> <url2> ... <urln>")
		os.Exit(1)
	}

	for _, url := range os.Args[1:] {
		wg.Add(1)
		go sendRequestAsync("https://" + url)
		//sendRequestSync("https://" + url)
	}

	wg.Wait()
}

func sendRequestAsync(url string) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mut.Lock()
	defer mut.Unlock()
	fmt.Printf("[%d] %s\n", res.StatusCode, url)
}

// func sendRequestSync(url string) {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	fmt.Printf("[%d] %s\n", res.StatusCode, url)
// }
