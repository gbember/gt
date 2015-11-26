// line_test.go
package navmesh

import (
	"testing"
)


func Benchmark_line_isIntersect(b *testing.B){
	l1 := line{sp: point{0, 0}, ep: point{10, 10}}
	l2 := line{sp: point{100, 100}, ep: point{5, 5}}
	for i:=0;i<b.N;i++{
		l1.isIntersectNotBothP(l2)
	}
}