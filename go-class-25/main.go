package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, url string, ch chan<- result) {
	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	// req now has url and context, meaning 1 second timeout is injected into http get
	// now we will get error from `http.DefaultClient.Do(req)` call rather than results
	if resp, err := http.DefaultClient.Do(req); err != nil {
		ch <- result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		ch <- result{url, nil, t}
		resp.Body.Close()
	}
}

func main() {
	// UPDATED GET CODE FROM go-class-24

	results := make(chan result)
	list := []string{
		"https://amazon.com", // usually takes > 1 sec
		"https://google.com",
		"https://nytimes.com",
		"https://wsl.com",
	}

	// creating context with 1 second timeout
	// this approach is cleaner than creating a select with a stopper
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for _, url := range list {
		go get(ctx, url, results)
	}

	for range list {
		r := <-results
		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
		} else {
			log.Printf("%-20s %s\n", r.url, r.latency)
		}
	}
}
