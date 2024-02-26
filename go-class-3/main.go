package main

import (
	"fmt"
	"os"
)

func main() {
	var sum float64
	var n int

	for {
		var val float64

		// address of operator is "&"
		// Fscanln will pass in value to pointer of val (&val)
		if _, err := fmt.Fscanln(os.Stdin, &val); err != nil {
			break
		}

		sum += val
		n++
	}

	if n == 0 {
		fmt.Fprintln(os.Stderr, "no values")
		os.Exit(-1)
	}

	fmt.Println("The average is", sum/float64(n))
}
