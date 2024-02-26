package main

import "testing"

// BENCHMARKING
// benchmarks live in test files ending with _test.go
// you run benchmarks with go test -bench; go only runs the `BenchmarkXXX` functions --> go test -bench=. ./fibonacci_test.go

func Fib(n int, recusive bool) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	default:
		if recusive {
			return Fib(n-1, true) + Fib(n-2, true) // recursive fib calcs is less efficient
		}
		a, b := 0, 1

		for i := 1; i < n; i++ {
			a, b = b, a+b
		}
		return b
	}
}

func BenchmarkFib20T(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(20, true)
	}
}

func BenchmarkFib20F(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(20, false)
	}
}
