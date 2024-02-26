package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const url = "https://jsonplaceholder.typicode.com"

type todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	// appending some args to url and getting response from server
	resp, err := http.Get(url + "/todos/1")

	// ORIGINAL:
	// If you run `go run ./client/ preeti` while running the server, we will
	// get the response from the server in ../server/ as specified by handler
	// function, which would be "Hello, world! from preeti"
	//resp, err := http.Get("http://localhost:8080/" + os.Args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// response body needs to be closed because otherwise it won't close the socket
	// and then we will run out. defer will close response at end of main
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// ORIGINAL:
		// fmt.Println(string(body))

		var item todo
		err = json.Unmarshal(body, &item)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		fmt.Printf("%#v\n", item)
	}
}
