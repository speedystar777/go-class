package main

import (
	"testing"
)

// for coverage run: go test -cover
// for more detailed coverage output: go test -coverprofile=c.out -covermode=count
// to visualize the output: go tool cover -html=c.out

var unknown = `{
	"id": 1,
	"name": "bob",
	"addr": {
		"street": "Lazy Lane",
		"city": "Exit",
		"zip": "9999"
	},
	"extra": 21.1
}`

// test checking to see that certain pieces of data are in an unknown object
func TestContains(t *testing.T) {
	var known = []string{
		`{"id": 1}`,
		`{"extra": 21.1}`,
		`{"name": "bob"}`,
		`{"addr": {"street": "Lazy Lane", "city": "Exit"}}`,
	}

	for _, k := range known {
		if err := CheckData([]byte(k), []byte(unknown)); err != nil {
			t.Errorf("invalid: %s (%s)\n", k, err)
		}
	}
}

func TestNotContains(t *testing.T) {
	var known = []string{
		`{"id": 2}`,
		`{"pid": 2}`,
		`{"name": "bobby"}`,
		`{"first": "bob"}`,
		`{"addr": {"street": "Lazy Lane", "city": "Alpha"}}`,
		`{"name": {"avenue": "Lazy Ave"}}`,
		`{"city": {"avenue": "Lazy Ave"}}`,
	}

	for _, k := range known {
		if err := CheckData([]byte(k), []byte(unknown)); err == nil {
			t.Errorf("false positive: %s\n", k)
		} else {
			t.Log(err)
		}
	}
}
