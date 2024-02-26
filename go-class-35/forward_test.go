package main

import (
	"math/rand"
	"testing"
)

// go test -bench=. ./forward_test.go

const defaultChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func randString(length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

type forwarder interface {
	forward(string) int
}

type thing1 struct {
	t forwarder
}

func (t1 *thing1) forward(s string) int {
	return t1.t.forward(s)
}

type thing2 struct {
	t forwarder
}

func (t2 *thing2) forward(s string) int {
	return t2.t.forward(s)
}

type thing3 struct {
}

func (t3 *thing3) forward(s string) int {
	return len(s)
}

func length(s string) int {
	return len(s)
}

// creates random string and finds length of it
func BenchmarkDirect(b *testing.B) {
	r := randString(rand.Intn(24), defaultChars)
	h := make([]int, b.N)

	b.ResetTimer() // not looking time to create string

	for i := 0; i < b.N; i++ {
		h[i] = length(r)
	}
}

func BenchmarkForward(b *testing.B) {
	r := randString(rand.Intn(24), defaultChars)
	h := make([]int, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ { // creating extra objects and forwarding it is more expensive
		var t3 forwarder = &thing3{}
		var t2 forwarder = &thing2{t3}
		var t1 forwarder = &thing1{t2}

		h[i] = t1.forward(r)
	}
}
