package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func matchNum(key string, exp float64, data map[string]interface{}) bool {
	if v, ok := data[key]; ok {
		if val, ok := v.(float64); ok && val == exp {
			return true
		}
	}
	return false
}

func matchString(key string, exp string, data map[string]interface{}) bool {
	if v, ok := data[key]; ok {
		if val, ok := v.(string); ok && strings.EqualFold(val, exp) {
			return true
		}
	}
	return false
}

func contains(exp, data map[string]interface{}) error {
	for k, v := range exp {
		switch x := v.(type) {
		case float64:
			if !matchNum(k, x, data) {
				return fmt.Errorf("%s unmatched (%d)", k, int(x))
			}
		case string:
			if !matchString(k, x, data) {
				return fmt.Errorf("%s unmatched (%s)", k, x)
			}
		case map[string]interface{}:
			if val, ok := data[k]; !ok {
				return fmt.Errorf("%s missing", k)
			} else if unk, ok := val.(map[string]interface{}); ok {
				if err := contains(x, unk); err != nil {
					return fmt.Errorf("%s unmatched in %#v: %s", k, x, err)
				}
			} else {
				return fmt.Errorf("%s wring in %#v", k, val)
			}
		}
	}
	return nil
}

func CheckData(want, got []byte) error {
	var w, g map[string]interface{}

	if err := json.Unmarshal(want, &w); err != nil {
		return err
	}

	if err := json.Unmarshal(got, &g); err != nil {
		return err
	}

	return contains(w, g)
}

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
