package main

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Page        int      `json:"page"`
	Words       []string `json:"words, omitempty"`
	randomfield int      `json:"randomfield"` // randomfield is a private field and therefore not exported in json even though there is no omitempty 
}

func main() {
	r := &Response{Page: 1, Words: []string{"one", "two", "three"}}
	j, _ := json.Marshal(r)
	fmt.Println(string(j))
	fmt.Printf("%#v\n", r)
	var r2 Response
	_ = json.Unmarshal(j, &r2)
	fmt.Printf("%#v\n", r2)
}
