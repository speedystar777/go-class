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
	// UPDATED GET CODE FROM go-class-23
	// we want to stop reading after 1 second
	stopper := time.After(1 * time.Second)

	results := make(chan result)
	list := []string{
		"https://amazon.com",
		"https://google.com",
		"https://nytimes.com",
		"https://wsl.com",
	}

	for _, url := range list {
		go get(url, results)
	}

	for range list {
		// reading from all channels
		select {
		case r := <-results:
			if r.err != nil {
				log.Printf("%-20s %s\n", r.url, r.err)
			} else {
				log.Printf("%-20s %s\n", r.url, r.latency)
			}
		// stopper is available to read after 1 second and stops program
		// result is we wont get results from every url (amazon.com usually takes > 1 sec)
		case t := <-stopper:
			log.Fatalf("timeout %s", t)
		}
	}
}
