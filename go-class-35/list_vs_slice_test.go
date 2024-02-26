package main

import "testing"

// go test -bench=. ./list_vs_slice_test.go
// go test -bench=. -benchmem ./list_vs_slice_test.go --> includes memory allocations

type node struct {
	v *int  // value in list
	t *node // pointers are less efficient to retrieve from cache and use more memory
}

func insert(i int, h *node) *node {
	t := &node{&i, nil}
	if h != nil {
		h.t = t
	}
	return t
}

func mkList(n int) *node {
	var h, t *node
	h = insert(0, h)
	t = insert(1, h)

	for i := 2; i < n; i++ {
		t = insert(i, t)
	}

	return h
}

func sumList(h *node) (i int) {
	for n := h; n != nil; n = n.t {
		i += *h.v
	}
	return
}

func mkSlice(n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

func sumSlice(l []int) (i int) {
	for _, v := range l {
		i += v
	}
	return
}

func BenchmarkList(b *testing.B) {
	l := mkList(1200)
	b.ResetTimer() // so that we are not including mkList/mkSlice in time
	for n := 0; n < b.N; n++ {
		sumList(l)
	}
}

func BenchmarkSlice(b *testing.B) { // efficiency of cache is better in slice
	l := mkSlice(1200)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sumSlice(l)
	}
}
