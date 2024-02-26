package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// run server with go run -race ./server
// then run this file in another terminal with go run .
type sku struct {
	item, price string
}

var items = []sku{
	{"shoes", "46"},
	{"socks", "6"},
	{"sandals", "27"},
	{"clogs", "36"},
	{"pants", "30"},
	{"shorts", "20"},
}

func doQuery(cmd, params string) error {
	resp, err := http.Get("http://localhost:8080/" + cmd + "?" + params)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err %s = %v\n", params, err)
		return err
	}

	defer resp.Body.Close()

	fmt.Fprintf(os.Stderr, "got %s = %d (no err)\n", params, resp.StatusCode)
	return nil
}

func runAdds() {
	for {
		for _, s := range items {
			if err := doQuery("create", "item="+s.item+"&price="+s.price); err != nil {
				return
			}
		}
	}
}

func runUpdates() {
	for {
		for _, s := range items {
			if err := doQuery("update", "item="+s.item+"&price="+s.price); err != nil {
				return
			}
		}
	}
}

func runDeletes() {
	for {
		for _, s := range items {
			if err := doQuery("delete", "item="+s.item); err != nil {
				return
			}
		}
	}
}

func main() {

	go runAdds()
	go runDeletes()
	go runUpdates()

	time.Sleep(10 * time.Second)
}
