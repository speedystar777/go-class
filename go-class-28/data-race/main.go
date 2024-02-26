package main

import (
	"fmt"
	"sync"
)

// will start 1000 goroutines
// all goroutines will try to increment variable n, which they can all see and modify
func do() int {
	var n int64
	var w sync.WaitGroup

	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			n++ // DATA RACE: read, modify, and write of a variable shared between goroutines
			w.Done()
		}()
	}

	w.Wait()      // waitGroup will make sure all 1000 goroutines will run
	return int(n) // however, output will not be 1000 due to data race
}

func main() {
	fmt.Println(do())
}
