package main

import (
	"log"
	"net/http"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

// when passing channel as func parameter, we can restrict use of it
// in this case we are only giving write permission to channel for this func
func get(url string, ch chan<- result) {
	start := time.Now()

	if resp, err := http.Get(url); err != nil {
		ch <- result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		ch <- result{url, nil, t}
		resp.Body.Close()
	}
}

func main() {
	// creating channel of result type
	results := make(chan result)
	list := []string{
		"https://amazon.com",
		"https://google.com",
		"https://nytimes.com",
		"https://wsl.com",
	}

	// starting a goroutine for each url in list
	for _, url := range list {
		go get(url, results)
	}
	// if there is no data to read in the channel and you go to read it,
	// you will have to wait for data -> you will be blocked
	//--------
	// if we were to range over channel itself, it would read until the channel
	// closes, but we are not closing the channel, because a channel can only be
	// closed once, so we must close after all the goroutines are finished, but
	// we don't know the order they will finish
	//--------
	// therefore it is best to range over lists since we know, based on the code
	// for the `get` function, each list has a result written to the channel
	// whether or not the fetch function errors
	for range list {
		// arrow before does read on results channel and puts it into var r
		r := <-results

		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
		} else {
			log.Printf("%-20s %s\n", r.url, r.latency)
		}
	}
}
