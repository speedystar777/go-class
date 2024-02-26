package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type response struct {
	Item   string `json:"item"`
	Album  string
	Title  string
	Artist string
}

type respWrapper struct {
	response
}

var j1 = `{
	"item": "album",
	"album": {"title": "Dark Side of the Moon"}
}`

var j2 = `{
	"item": "song",
	"song": {"title": "Bella Donna", "artist": "Stevie Nicks"}
}`

// custom unmarshal method will be called since respWrapper implements interface that requires UnmarshalJSON
func (r *respWrapper) UnmarshalJSON(b []byte) (err error) {
	var raw map[string]interface{}

	_ = json.Unmarshal(b, &r.response)
	err = json.Unmarshal(b, &raw)

	switch r.Item {
	case "album":
		inner, ok := raw["album"].(map[string]interface{}) // using type assertions to weave through json map
		if ok {
			if album, ok := inner["title"].(string); ok {
				r.Album = album
			}
		}
	case "song":
		inner, ok := raw["song"].(map[string]interface{})
		if ok {
			if title, ok := inner["title"].(string); ok {
				r.Title = title
			}
			if artist, ok := inner["artist"].(string); ok {
				r.Artist = artist
			}
		}
	}

	return err
}

func main() {
	var resp1, resp2 respWrapper
	var err error

	if err = json.Unmarshal([]byte(j1), &resp1); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", resp1.response)

	if err = json.Unmarshal([]byte(j2), &resp2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", resp2.response)
}
