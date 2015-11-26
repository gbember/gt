// line.go
package navmesh

type line struct {
	sp point
	ep point
}

//检测是否与线相交
func (l line) isIntersect(ol line) bool {
	lmaxx, lminx := maxmin(l.sp.x, l.ep.x)
	olmaxx, olminx := maxmin(ol.sp.x, ol.ep.x)
	lmaxy, lminy := maxmin(l.sp.y, l.ep.y)
	olmaxy, olminy := maxmin(ol.sp.y, ol.ep.y)
	if !(lmaxx >= olminx && olmaxx >= lminx &&
		lmaxy >= olminy && olmaxy >= lminy) {
		return false
	}
	z1 := (ol.sp.x-l.sp.x)*(l.ep.y-l.sp.y) - (ol.sp.y-l.sp.y)*(l.ep.x-l.sp.x)
	z2 := (ol.ep.x-l.sp.x)*(l.ep.y-l.sp.y) - (ol.ep.y-l.sp.y)*(l.ep.x-l.sp.x)
	z3 := (l.sp.x-ol.sp.x)*(ol.ep.y-ol.sp.y) - (l.sp.y-ol.sp.y)*(ol.ep.x-ol.sp.x)
	z4 := (l.ep.x-ol.sp.x)*(ol.ep.y-ol.sp.y) - (l.ep.y-ol.sp.y)*(ol.ep.x-ol.sp.x)
	return z1*z2 <= 0 && z3*z4 <= 0
}
