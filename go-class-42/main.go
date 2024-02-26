package main

import (
	"fmt"
	"strconv"
	"strings"
)

// create vector of type T, which can be anything (any valid type; no constraint on T)
type Vector[T any] []T

// generic push function
// takes pointer so input is modified
func (v *Vector[T]) Push(x T) {
	*v = append(*v, x)
}

func Map[F, T any](s []F, f func(F) T) []T {
	r := make([]T, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// can't create String method for standard types like int, so we create a type
type num int

func (n num) String() string {
	return strconv.Itoa(int(n))
}

// here, T must be a Stringer, i.e., implement that String() method, so we create type num, which is an int that has String() method
type StringableVector[T fmt.Stringer] []T

// string method on StringableVector for formatting
func (s StringableVector[T]) String() string {
	var sb strings.Builder

	sb.WriteString("<<")

	for i, v := range s {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString(">>")

	return sb.String()
}

func main() {
	s := Vector[int]{}

	s.Push(1)
	s.Push(2)

	// both map functions will work
	t1 := Map(s, strconv.Itoa) // map function did not need type parameter specifications, i.e., Map[int, string](...)
	// this is because the compiler has enough information based on input parameters to determine types
	t2 := Map([]int{1, 2, 3}, strconv.Itoa)

	fmt.Println(s)
	fmt.Printf("%#v\n", t1)
	fmt.Printf("%#v\n", t2)

	var sv StringableVector[num] = []num{1, 2, 3} // must have num in square brackets after StringableVector
	// the fact that sv is assigned to a string of nums is not enough information for compiler to infer that sv is StringableVector of type num
	fmt.Println(sv)
}
