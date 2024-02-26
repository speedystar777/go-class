package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, url string, ch chan<- result) {
	var r result
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second).C
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	// req now has url and context, meaning 1 second timeout is injected into http get
	// now we will get error from `http.DefaultClient.Do(req)` call rather than results
	if resp, err := http.DefaultClient.Do(req); err != nil {
		r = result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		r = result{url, nil, t}
		resp.Body.Close()
	}

	for {
		select {
		// if r can be sent on channel, do that and return from function; end of goroutine
		case ch <- r:
			return
		// if we are blocked because nobody is ready to receive, then hit this ticker
		// if ch has buffer, then they will not tick
		case <-ticker:
			log.Println("tick", r)
		}
	}
}

func first(ctx context.Context, urls []string) (*result, error) {
	// BUFFER TO AVOID LEAKING AS BELOW: sender can send, can receiver can read later
	// in our case, we will not read later since we only want to read first result, but buffer breaks up the sending/receiving
	// with buffer there is no ticking
	// results := make(chan result, len(urls))
	results := make(chan result) // this results channel has no buffer
	// for an unbuffered channel, if someone wants to send, someone else has to be able to receive

	ctx, cancel := context.WithCancel(ctx)

	// need to cancel even if timeout doesn't happen to release resources
	// http operation is being canceled, not go routines
	defer cancel()

	for _, url := range urls {
		// starting go routines (parallel gets)
		go get(ctx, url, results) // after we cancel, other go routines that want to write to channel will get blocked
	}

	select {
	// normal case, get response and return it; this return will cause deferred cancellation, so all other requests will get cancelled
	// we will not be listening to channel after first result is read
	// the only reader is gone so no other go routines can write to this channel --> leaky goroutines cause memory leak
	case r := <-results:
		return &r, nil
	// handles the case of the context becoming done
	// our cancellation will not cause this select case to happen, this case will happen if parent context times out
	// because we are being given context, which could have parents,
	// and have select in this method we need to be prepared for context being cancelled from above
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func main() {
	list := []string{
		"https://amazon.com", // usually takes > 1 sec
		"https://google.com",
		"https://nytimes.com",
		"https://wsl.com",
	}

	// this function is returning the first response
	r, _ := first(context.Background(), list)

	if r.err != nil {
		log.Printf("%-20s %s\n", r.url, r.err)
	} else {
		log.Printf("%-20s %s\n", r.url, r.latency)
	}

	time.Sleep(9 * time.Second)
	log.Println("quit anyway...", runtime.NumGoroutine(), "still running")

}
