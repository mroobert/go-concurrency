package main

import (
	"fmt"
	"time"
)


func DoneMain(){
	done := make(chan struct{})

	go doWork(done)
	time.Sleep(5 * time.Second)
	close(done)
}

func doWork(done <-chan struct{}){
	for {
		select{
		case <-done:
			fmt.Println("Done!")
			return
		default:
			fmt.Println("Doing work...")
		}
	}
}