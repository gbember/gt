// line_test.go
package navmesh

import (
	"testing"
)

func Benchmark_line_isIntersect(b *testing.B) {
	l1 := line{sp: Point{0, 0}, ep: Point{10, 10}}
	l2 := line{sp: Point{100, 100}, ep: Point{5, 5}}
	for i := 0; i < b.N; i++ {
		l1.isCross(l2)
	}
}
