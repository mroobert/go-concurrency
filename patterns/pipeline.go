package main

import "fmt"


func PipelineMain(){
	done := make(chan struct{})
	defer close(done)

	values := []int{1,2,3,4}

	// 2,4,6,8 -> 3, 5, 7, 9 -> 6, 10, 14, 18
	// pipeline := Multiply(done,Add(done,Multiply(done,Generator(done,values),2),1),2)

	// 1,2,3,4 -> 2,3,4,5 -> 4,6,8,10
	pipeline := Multiply(done,Add(done,Generator(done,values),1),2)
	
	for v := range pipeline{
		fmt.Println(v)
	}
}


// Generator converts a discrete set of int values into a stream of data on a channel.
func Generator(done <-chan struct{}, values []int) <-chan int{
	out := make(chan int)
	go func(){
		defer close(out)
		for _, v := range values {
			select{
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
} 

// Add uses the adder value to increase each value from the channel.
func Add(done <-chan struct{}, values <-chan int, adder int) <-chan int{
	out := make(chan int)
	go func(){
		defer close(out)
		for v := range values{
			select{
			case <- done:
				return
			case out <- v + adder:
			}
	}
	}()

	return out
}

// Multiply uses the multiplier value to multiply each value from the channel.
func Multiply(done <-chan struct{}, values <-chan int, multiplier int) <-chan int{
	out := make(chan int)
	go func(){
		defer close(out)
		for v := range values{
			select{
			case <-done:
				return
			case out <- v * multiplier:
			}
		}
	 }()

	return out
}