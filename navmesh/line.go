// line.go
package navmesh

import (
	"math"
)

type line struct {
	sp Point
	ep Point
}

//检测是否与线相交(交叉点包括线的起点和终点)
func (l line) isIntersect(ol line) bool {
	lmaxx, lminx := maxmin(l.sp.X, l.ep.X)
	olmaxx, olminx := maxmin(ol.sp.X, ol.ep.X)

	if !(lmaxx >= olminx && olmaxx >= lminx) {
		return false
	}

	lmaxy, lminy := maxmin(l.sp.Y, l.ep.Y)
	olmaxy, olminy := maxmin(ol.sp.Y, ol.ep.Y)

	if !(lmaxy >= olminy && olmaxy >= lminy) {
		return false
	}

	z1 := (ol.sp.X-l.sp.X)*(l.ep.Y-l.sp.Y) - (ol.sp.Y-l.sp.Y)*(l.ep.X-l.sp.X)
	z2 := (ol.ep.X-l.sp.X)*(l.ep.Y-l.sp.Y) - (ol.ep.Y-l.sp.Y)*(l.ep.X-l.sp.X)
	if z1*z2 > 0 {
		return false
	}
	z3 := (l.sp.X-ol.sp.X)*(ol.ep.Y-ol.sp.Y) - (l.sp.Y-ol.sp.Y)*(ol.ep.X-ol.sp.X)
	z4 := (l.ep.X-ol.sp.X)*(ol.ep.Y-ol.sp.Y) - (l.ep.Y-ol.sp.Y)*(ol.ep.X-ol.sp.X)
	if z3*z4 > 0 {
		return false
	}
	return true
}

//检测是否与线互相穿过(只有一个交点且交点不是起点和终点)
func (l line) isCross(ol line) bool {
	lmaxx, lminx := maxmin(l.sp.X, l.ep.X)
	olmaxx, olminx := maxmin(ol.sp.X, ol.ep.X)

	if !(lmaxx >= olminx && olmaxx >= lminx) {
		return false
	}

	lmaxy, lminy := maxmin(l.sp.Y, l.ep.Y)
	olmaxy, olminy := maxmin(ol.sp.Y, ol.ep.Y)

	if !(lmaxy >= olminy && olmaxy >= lminy) {
		return false
	}

	z1 := (ol.sp.X-l.sp.X)*(l.ep.Y-l.sp.Y) - (ol.sp.Y-l.sp.Y)*(l.ep.X-l.sp.X)
	z2 := (ol.ep.X-l.sp.X)*(l.ep.Y-l.sp.Y) - (ol.ep.Y-l.sp.Y)*(l.ep.X-l.sp.X)
	if z1*z2 >= 0 {
		return false
	}
	z3 := (l.sp.X-ol.sp.X)*(ol.ep.Y-ol.sp.Y) - (l.sp.Y-ol.sp.Y)*(ol.ep.X-ol.sp.X)
	z4 := (l.ep.X-ol.sp.X)*(ol.ep.Y-ol.sp.Y) - (l.ep.Y-ol.sp.Y)*(ol.ep.X-ol.sp.X)
	if z3*z4 >= 0 {
		return false
	}
	return true
}

//线是否与多边形的边相交(交叉点包括线的起点和终点)
func (l line) isIntersectConvexPolygon(nm *NavMesh, cp *convexPolygon) bool {
	length := len(cp.pindexs) - 1
	ol := line{sp: nm.points[cp.pindexs[length-1]], ep: nm.points[cp.pindexs[0]]}
	if l.isIntersect(ol) {
		return true
	}
	for i := 1; i < length; i++ {
		ol.sp = ol.ep
		ol.ep = nm.points[cp.pindexs[i]]
		if l.isIntersect(ol) {
			return true
		}
	}

	return false
}

//线是否穿过多边形的边(各顶点和边不算)
func (l line) isCrossConvexPolygon(nm *NavMesh, cp *convexPolygon) bool {
	length := len(cp.pindexs) - 1
	ol := line{sp: nm.points[cp.pindexs[length-1]], ep: nm.points[cp.pindexs[0]]}
	if l.isCross(ol) {
		return true
	}
	for i := 1; i < length; i++ {
		ol.sp = ol.ep
		ol.ep = nm.points[cp.pindexs[i]]
		if l.isCross(ol) {
			return true
		}
	}
	return false
}

//线是否穿过多边形不可穿过的边(顶点不算、与不可穿过边交叉不)
func (l line) isCrossConvexPolygoNoPassLines(nm *NavMesh, cp *convexPolygon) bool {
	for _, l1 := range cp.lines {
		if l.isCross(l1) {
			return true
		}
	}
	return false
}

//得到线段穿过的格子id列表
func (l line) getAcossGridNums(gsize int64, maxVNum int64) []int64 {
	gnum1 := l.sp.getGridNum(gsize, maxVNum)
	gnum2 := l.ep.getGridNum(gsize, maxVNum)
	if gnum1 == gnum2 {
		return []int64{gnum1}
	}
	gidList := make([]int64, 0, 20)
	//在同一行
	if int64(math.Abs(float64(gnum1-gnum2))) < gsize {
		if gnum1 > gnum2 {
			for ; gnum2 <= gnum1; gnum2++ {
				gidList = append(gidList, gnum2)
			}
		} else {
			for ; gnum1 <= gnum2; gnum1++ {
				gidList = append(gidList, gnum1)
			}
		}
		return gidList
	}
	//在同一列
	if gnum1%maxVNum == gnum2%maxVNum {
		if gnum1 > gnum2 {
			for ; gnum2 <= gnum1; gnum2 += maxVNum {
				gidList = append(gidList, gnum2)
			}
		} else {
			for ; gnum1 <= gnum2; gnum1 += maxVNum {
				gidList = append(gidList, gnum1)
			}
		}
		return gidList
	}
	x := l.ep.X - l.sp.X
	y := l.ep.Y - l.sp.Y
	tan := y / x
	a := l.ep.Y - tan*l.ep.X
	gid := l.sp.getGridNum(gsize, maxVNum)
	gidList = append(gidList, gid)
	p := Point{}
	if x > 0 {
		max := l.ep.X / gsize * gsize
		x = l.sp.X/gsize*gsize + gsize
		for ; x <= max; x += gsize {
			y = tan*x + a
			p.X, p.Y = x, y
			gid = p.getGridNum(gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	} else {
		min := l.ep.X / gsize * gsize
		x = l.sp.X/gsize*gsize - gsize
		for ; x >= min; x -= gsize {
			y = tan*x + a
			p.X, p.Y = x, y
			gid = p.getGridNum(gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	}
	if l.ep.Y-l.sp.Y > 0 {
		max := l.ep.Y / gsize * gsize
		y = l.sp.Y/gsize*gsize + gsize
		for ; y <= max; y += gsize {
			x = (y - a) / tan
			p.X, p.Y = x, y
			gid = p.getGridNum(gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	} else {
		min := l.ep.Y / gsize * gsize
		y = l.sp.Y/gsize*gsize - gsize
		for ; y >= min; y -= gsize {
			x = (y - a) / tan
			p.X, p.Y = x, y
			gid = p.getGridNum(gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	}
	return gidList
}

//长度没有开方
func (l line) Distance2() int64 {
	x := l.ep.X - l.sp.X
	y := l.ep.Y - l.sp.Y
	return x*x + y*y
}

//长度已开方
func (l line) Distance() float64 {
	x := l.ep.X - l.sp.X
	y := l.ep.Y - l.sp.Y
	return math.Sqrt(float64(x*x + y*y))
}

//是否是同一条线段(起点等于起点终点等于终点或者起点等于终点终点等于起点)
func (l line) isSame(l1 line) bool {
	return (l.sp == l1.sp && l.ep == l1.ep) || (l.ep == l1.sp && l.sp == l1.ep)
}
