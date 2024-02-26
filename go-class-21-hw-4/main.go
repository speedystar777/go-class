package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	if _, ok := db[item]; ok {
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

	db[item] = dollars(p)
	fmt.Fprintf(w, "added %s with price %s", item, db[item])
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	if _, ok := db[item]; !ok {
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

	db[item] = dollars(p)
	fmt.Fprintf(w, "updated %s with new price %s", item, db[item])
}

func (db database) fetch(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	fmt.Fprintf(w, "item %s has price %s", item, db[item])
}

func (db database) drop(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	delete(db, item)
	fmt.Fprintf(w, "dropped item %s", item)
}

func main() {
	db := database{
		"shoes": 50,
		"socks": 5,
	}

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/fetch", db.fetch)
	http.HandleFunc("/drop", db.drop)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
