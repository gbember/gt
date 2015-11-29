package navmesh

import (
	"testing"
)

var (
	s1 []bool
	s2 []bool
)

func init() {
	s1 = make([]bool, 200, 200)
	s2 = make([]bool, 200, 200)
}

func BenchmarkCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		copy(s2, s1)
	}
}

func BenchmarkForEQ(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s2); j++ {
			s2[j] = false
		}
	}
}
