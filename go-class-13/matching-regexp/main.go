package main

import (
	"fmt"
	"regexp"
	"strings"
)

// find all string with one or more of lowercase character b
// can also replace those string
func main() {
	te := "aba abba abbba"

	re := regexp.MustCompile(`b+`)
	mm := re.FindAllString(te, -1)
	id := re.FindAllStringIndex(te, -1)

	fmt.Println(mm) // [b bb bbb]
	fmt.Println(id) // [[1 2] [5 7] [10 13]]

	for _, d := range id {
		fmt.Println(te[d[0]:d[1]])
	}

	up := re.ReplaceAllStringFunc(te, strings.ToUpper)
	fmt.Println(up)
}
