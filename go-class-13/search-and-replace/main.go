package main

import (
	"fmt"
	"regexp"
)

// parentheses indicate capture groups; phone number gets broken into 3 parts
var ph = regexp.MustCompile(`\(([[:digit:]]{3})\) ([[:digit:]]{3})-([[:digit:]]{4})`)

func main() {
	orig := "(214) 514-9548"
	// first element is entire string that matches, following elements are possible submatches based on capture groups
	match := ph.FindStringSubmatch(orig)

	fmt.Printf("Number 1: %q\n", match)
	if len(match) > 3 {
		fmt.Printf("+1 %s-%s-%s\n", match[1], match[2], match[3])
	}

	orig2 := "call me at (214) 514-9548 today"
	match2 := ph.FindStringSubmatch(orig2)
	fmt.Printf("Number 2: %q\n", match2)                    // same output as above
	intl := ph.ReplaceAllString(orig2, "+1 ${1}-${2}-${3}") // ${1} indicates submatch at index 1, etc.
	fmt.Println(intl)
}
