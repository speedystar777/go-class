package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// run this with race detector condition turned on: go run -race ./server
type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	mu sync.Mutex // add mutex to prevent data race conditions
	db map[string]dollars
}

func (d *database) list(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for item, price := range d.db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (d *database) create(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	if _, ok := d.db[item]; ok {
		msg := fmt.Sprintf("duplicate item: %q", item)
		http.Error(w, msg, http.StatusBadRequest) // 400
		return                                    // if there is an error, handle error and get out
		// without return you would signal error and then try to do regular processing
	}
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) // 400
		return
	}

	d.db[item] = dollars(p)
	fmt.Fprintf(w, "added %s with price %s", item, d.db[item])
}

func (d *database) update(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	if _, ok := d.db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) // 400
		return
	}

	d.db[item] = dollars(p)
	fmt.Fprintf(w, "updated %s with new price %s", item, d.db[item])
}

func (d *database) fetch(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")

	if _, ok := d.db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	fmt.Fprintf(w, "item %s has price %s", item, d.db[item])
}

func (d *database) drop(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")

	if _, ok := d.db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	delete(d.db, item)
	fmt.Fprintf(w, "dropped item %s", item)
}

func main() {
	db := database{
		db: map[string]dollars{
			"shoes": 50,
			"socks": 5,
		},
	}

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/fetch", db.fetch)
	http.HandleFunc("/drop", db.drop)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
