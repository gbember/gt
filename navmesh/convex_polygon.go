// convex_polygon.go
package navmesh

//凸多边形
type ConvexPolygon struct {
	//多边形的所有顶点(顺时针排序)
	ps []point
}

//是否包含点
func (cp *ConvexPolygon)isContainPoint(p point)bool{
	return false
}

//是否包含格子
func (cp *ConvexPolygon)isContainGrid(gridNum int64)bool{
	return false
}

