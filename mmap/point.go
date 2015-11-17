// point.go
package mmap

type point struct {
	x int32
	y int32
}

//点向量相减
func (p point) sub(p1 point) point {
	return point{p.x - p1.x, p.y - p1.y}
}

//点向量叉乘
func (p point) CrossProduct(p1 point) int32 {
	return p.x*p1.y - p.y*p1.x
}
