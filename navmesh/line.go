// line.go
package navmesh

import(
	"math"
)

type line struct {
	sp point
	ep point
}

//检测是否与线相交(交叉点包括线的起点和终点)
func (l line) isIntersect(ol line) bool {
	lmaxx, lminx := maxmin(l.sp.x, l.ep.x)
	olmaxx, olminx := maxmin(ol.sp.x, ol.ep.x)
	
	if !(lmaxx >= olminx && olmaxx >= lminx) {
		return false
	}
	
	lmaxy, lminy := maxmin(l.sp.y, l.ep.y)
	olmaxy, olminy := maxmin(ol.sp.y, ol.ep.y)
	
	if !(lmaxy >= olminy && olmaxy >= lminy){
		return false
	}
	
	z1 := (ol.sp.x-l.sp.x)*(l.ep.y-l.sp.y) - (ol.sp.y-l.sp.y)*(l.ep.x-l.sp.x)
	z2 := (ol.ep.x-l.sp.x)*(l.ep.y-l.sp.y) - (ol.ep.y-l.sp.y)*(l.ep.x-l.sp.x)
	if z1*z2 > 0 {
		return false
	}
	z3 := (l.sp.x-ol.sp.x)*(ol.ep.y-ol.sp.y) - (l.sp.y-ol.sp.y)*(ol.ep.x-ol.sp.x)
	z4 := (l.ep.x-ol.sp.x)*(ol.ep.y-ol.sp.y) - (l.ep.y-ol.sp.y)*(ol.ep.x-ol.sp.x)
	if z3*z4 > 0 {
		return false
	}
	return true
}
//检测是否与线互相穿过(只有一个交点且交点不是起点和终点)
func (l line) isCross(ol line) bool {
	lmaxx, lminx := maxmin(l.sp.x, l.ep.x)
	olmaxx, olminx := maxmin(ol.sp.x, ol.ep.x)
	
	if !(lmaxx >= olminx && olmaxx >= lminx) {
		return false
	}
	
	lmaxy, lminy := maxmin(l.sp.y, l.ep.y)
	olmaxy, olminy := maxmin(ol.sp.y, ol.ep.y)
	
	if !(lmaxy >= olminy && olmaxy >= lminy){
		return false
	}
	
	z1 := (ol.sp.x-l.sp.x)*(l.ep.y-l.sp.y) - (ol.sp.y-l.sp.y)*(l.ep.x-l.sp.x)
	z2 := (ol.ep.x-l.sp.x)*(l.ep.y-l.sp.y) - (ol.ep.y-l.sp.y)*(l.ep.x-l.sp.x)
	if z1*z2 >= 0 {
		return false
	}
	z3 := (l.sp.x-ol.sp.x)*(ol.ep.y-ol.sp.y) - (l.sp.y-ol.sp.y)*(ol.ep.x-ol.sp.x)
	z4 := (l.ep.x-ol.sp.x)*(ol.ep.y-ol.sp.y) - (l.ep.y-ol.sp.y)*(ol.ep.x-ol.sp.x)
	if z3*z4 >= 0 {
		return false
	}
	return true
}

//线是否穿过多边形区域(各顶点和边不算)
func (l line)isCrossConvexPolygon(cp *ConvexPolygon)bool{
	length := len(cp.ps)-1
	ol := line{sp:cp.ps[length-1],ep:cp.ps[0]}
	if l.isCross(ol){
		return true
	}
	for i:=1;i<length;i++{
		ol.sp = ol.ep
		ol.ep = cp.ps[i]
		if l.isCross(ol){
			return true
		}
	}
	return false
}

//长度没有开方
func (l line) Distance2() int64 {
	x := l.ep.x - l.sp.x
	y := l.ep.y - l.sp.y
	return x*x + y*y
}
//长度已开方
func (l line) Distance() float64 {
	x := l.ep.x - l.sp.x
	y := l.ep.y - l.sp.y
	return math.Sqrt(float64(x*x + y*y))
}

//是否是同一条线段(起点等于起点终点等于终点或者起点等于终点终点等于起点)
func (l line) isSame(l1 line) bool {
	return (l.sp == l1.sp && l.ep == l1.ep) || (l.ep == l1.sp && l.sp == l1.ep)
}