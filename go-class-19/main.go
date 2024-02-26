package main

import (
	"fmt"
	"sort"
)

type Organ struct {
	Name   string
	Weight int
}

type Organs []Organ

func (s Organs) Len() int      { return len(s) }
func (s Organs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Methods Len and Swap will be promoted to ByName and ByWeight
type ByName struct{ Organs }
type ByWeight struct{ Organs }

func (s ByName) Less(i, j int) bool {
	return s.Organs[i].Name < s.Organs[j].Name
}

func (s ByWeight) Less(i, j int) bool {
	return s.Organs[i].Weight < s.Organs[j].Weight
}

func main() {
	s := []Organ{{"brain", 1340}, {"liver", 1494}, {"spleen", 162}, {"pancreas", 131}, {"heart", 290}}
	fmt.Println(s)
	// ByName and ByWeight are compatible with Sort interface since they have all three methods defined
	// Len and Swap, which are promoted and Less, which is defined separate for each type
	sort.Sort(ByWeight{s})             // converting type Organs to type ByWeight
	fmt.Println("Sort By Weight: ", s) // sort function is inline so it need to be printed before it is sorted again
	sort.Sort(ByName{s})               // converting type Organs to type ByName
	fmt.Println("Sort By Name: ", s)
	sort.Sort(sort.Reverse(ByName{s})) // convert input to reverse type first
	fmt.Println("Sort Reverse By Name: ", s)
}

// Can sort in reverse by using `sort.Reverse` which is defined as
// reverse is a private interface, not exposed to the user
// type reverse struct {
// 	// This embedded Interface permits Reverse to use the
// 	// methods of another Interface implementation
// 	Interface
// }
// // Less returns the opposite of the embedded implementation's Less method
// func (r reverse) Less(i, j int) bool {
// 	return r.Interface.Less(j, i)
// }
// // Reverse returns the reverse order for data
// func Reverse(data Interface) Interface {
// 	return &reverse{data}
// }
