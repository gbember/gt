// line_test.go
package navmesh

import (
	"testing"
)

func TestLine(t *testing.T) {
	l1 := line{sp: point{0, 0}, ep: point{10, 10}}
	l2 := line{sp: point{100, 100}, ep: point{20, 5}}
	t.Log(l1.isIntersect(l2))
}
