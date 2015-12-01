package test

import (
	"testing"
)

func Benchmark_1(b *testing.B) {
	var v []byte
	for i := 0; i < b.N; i++ {
		var x [8]byte
		v = x[:]
	}
	_ = v
}

func Benchmark_2(b *testing.B) {
	var v []byte
	for i := 0; i < b.N; i++ {
		var x = make([]byte, 8)
		v = x[:]
	}
	_ = v
}
