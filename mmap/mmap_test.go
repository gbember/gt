// mmap_test.go
package mmap

import "testing"

func TestMMap(t *testing.T) {
	//	Test()

	ps := []point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}

	t.Log(tt(point{5, 5}, ps))
	t.Log(tt(point{20, 5}, ps))
	t.Log(tt(point{10, 0}, ps))

}

func tt(p point, ps []point) bool {
	l := len(ps)
	for i := 0; i < l; i++ {
		if mul(p, ps[i], ps[(i+1)%l]) < 0 {
			return false
		}
	}
	return true
}

func mul(p point, p1 point, p2 point) int32 {
	return (p2.x-p1.x)*(p.y-p1.y) - (p.x-p1.x)*(p2.y-p1.y)
}
