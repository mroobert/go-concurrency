package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)


func SemaphoreMain(){
	WorkLoadConcurrentSemaphoreV2()
}


type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Workload represents a clasic use case where you request a resource from a server
// and decode the response into a struct.
//
// In this case we request 100 Tasks by doing 100 HTTP requests to the server and print the title of each one.
// For the moment the workload is not concurrent, but we will make it concurrent.
func WorkLoad(){
	var t Task

	for i:=0;i<100;i++{
		res, err := http.Get(fmt.Sprintf("http://jsonplaceholder.typicode.com/todos/%d", i))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			log.Fatal(err)
		}
		log.Println(t.Title)
	}
}

// WorkLoadConcurrent is the same as WorkLoad but it is concurrent.
// We are using waitGroup to tell the main goroutine to wait for all the goroutines to finish.
// In this way we can see the logs printed in the terminal.
func WorkLoadConcurrent(){
	var(
		t Task
		wg sync.WaitGroup
	)
	
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()

			res, err := http.Get(fmt.Sprintf("http://jsonplaceholder.typicode.com/todos/%d", i))
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
				log.Fatal(err)
			}
			log.Println(t.Title)
		}(i)
	}

	wg.Wait()
}

// WorkLoadConcurrentSemaphoreV1 is the same as WorkLoadConcurrent but it uses a semaphore 
// to limit the number of concurrent goroutines.
func WorkLoadConcurrentSemaphoreV1(){
	type token struct{}

	var(
		t Task
		wg sync.WaitGroup
	)

	sem := make(chan token, 10)

	wg.Add(100)
	for i := 0; i < 100; i++ {
		sem <- token{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			res, err := http.Get(fmt.Sprintf("http://jsonplaceholder.typicode.com/todos/%d", i))
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
				log.Fatal(err)
			}
			log.Println(t.Title)
		}(i)
	}

	wg.Wait()
	close(sem)
} 

// WorkLoadConcurrentSemaphoreV2 is the same as WorkLoadConcurrentSemaphoreV1 but it uses errgroup.
func WorkLoadConcurrentSemaphoreV2(){
	var t Task

	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(10)

	for i := 0; i < 100; i++ {
		i := i
		group.Go(func() error {
			res, err := http.Get(fmt.Sprintf("http://jsonplaceholder.typicode.com/todos/%d", i))
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
				return err
			}
			
			log.Println(t.Title)
			return nil
		})
	}

	err := group.Wait(); if err != nil {
		log.Fatal(err)
	}
}