package main

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

func FanMain() {

}

func fanIn(done <-chan any, channels ...<-chan any) <-chan any {
	var wg sync.WaitGroup
	multiplexedStream := make(chan any)

	multiplex := func(c <-chan any) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	// Select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	// Wait for all the reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

type City struct {
	Name       string
	Country    string
	Population int
}

// GenerateCities generates a stream of cities from a CSV file.
func GenerateCities(r io.Reader) chan City {
	out := make(chan City)
	go func() {
		defer close(out)

		reader := csv.NewReader(r)
		// discard first row, containg column names
		_, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			populationInt, err := strconv.Atoi(row[9])
			if err != nil {
				continue
			}

			out <- City{
				Name:       row[1],
				Country:    row[4],
				Population: populationInt,
			}
		}

	}()

	return out
}

// UpperCityName converts the city names to upper case.
func UpperCityName(rows <-chan City) <-chan City {
	out := make(chan City)
	go func() {
		defer close(out)

		for row := range rows {
			row.Name = strings.ToUpper(row.Name)
			out <- row
		}
	}()

	return out
}

// UpperCountryName converts the country names to upper case.
func UpperCountryName(rows <-chan City) <-chan City {
	out := make(chan City)
	go func() {
		defer close(out)

		for row := range rows {
			row.Country = strings.ToUpper(row.Country)
			out <- row
		}
	}()

	return out
}
