package main

import (
	"fmt"
	"log"
	"net/http"
)

// // DESIGN WITH `nextID` AS GLOBAL VARIABLE
// var nextID = make(chan int)

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "<h1> You got %d</h1>", <-nextID) // reader can't read unless there's somebody ready to write

// 	// nextID++ // UNSAFE: nextID is int
// 	// increment operation is read, modify, write
// 	// webserver built into Go standard library is concurrent so if you visit
// 	// webpage fast enough, there would be issues and numbers would be skipped
// }

// // function will run until server stops
// func counter() {
// 	for i := 0; ; i++ {
// 		nextID <- i // nothing can happen here unless someone is ready to read from it
// 	}
// }

// func main() {
// 	// will not result in explosive counting because function will pause at read and write steps
// 	go counter()
// 	http.HandleFunc("/", handler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// DESIGN WITH `nextID` AS INT CHANNEL
type nextCh chan int

func (ch nextCh) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> You got %d</h1>", <-ch)
}

func counter(ch chan<- int) {
	for i := 0; ; i++ {
		ch <- i
	}
}

func main() {
	var nextID nextCh = make(chan int) // need to specify nextCh type so it has handler function
	go counter(nextID)
	http.HandleFunc("/increment", nextID.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
