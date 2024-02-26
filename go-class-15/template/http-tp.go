package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type todo struct {
	UserID    int    `json:"userID"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var form = `
<h1>Todo #{{.ID}}</h1>
<div>{{printf "User %d" .UserID}}</div>
<div>{{printf "%s (completed: %t)" .Title .Completed}}</div>`

// Here webserver is acting as a client
func handler(w http.ResponseWriter, r *http.Request) {
	const base = "https://jsonplaceholder.typicode.com/"
	resp, err := http.Get(base + r.URL.Path[1:])

	if err != nil {
		// handler will write error to http response if base url is unavailable
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		// must return, otherwise handler logic will keep running
		return
	}

	// response body needs to be closed because otherwise it won't close the
	// socket and then we will run out. defer will close response at end of main
	defer resp.Body.Close()

	var item todo
	if err = json.NewDecoder(resp.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.New("mine")
	tmpl.Parse(form)
	tmpl.Execute(w, item)
}

// run `go run ./template/` and then append todos/1 to URL
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
