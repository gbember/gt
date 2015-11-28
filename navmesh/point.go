// vector2d.go
package navmesh

//像素点(向量)
type Point struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

//得到点所在的格子编号
func (p Point) getGridNum(gsize int64, maxVNum int64) int64 {
	return p.X/gsize + 1 + p.Y/gsize*maxVNum
}
