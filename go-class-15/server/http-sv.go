package main

import (
	"fmt"
	"log"
	"net/http"
)

const url = "https://jsonplaceholder.typicode.com"

type todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// This handler function will the write the string to the reponse writer,
// which is an input to Fpintf (which take any arg that can be written to).
// The response writer is handling an http response through a socket
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! from %s\n", r.URL.Path[1:])
}

// Once you run `go run ./server/` you will see the text above with whatever
// is after the slash in the url path. The server is receiving http request
// with URL and writes response
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
