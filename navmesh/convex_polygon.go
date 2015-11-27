// convex_polygon.go
package navmesh

import "errors"

//凸多边形
type convexPolygon struct {
	id    int        //区域编号
	ps    []point    //多边形的所有顶点(逆时针排序)
	lines []line     //该多边形不能穿过的线
	lcs   []*line2CP //相邻区域
}

//多边形相邻关系
type line2CP struct {
	l  line           //线
	cp *convexPolygon //区域
}

func newConvexPolygon(id int) *convexPolygon {
	cp := new(convexPolygon)
	cp.lines = make([]line, 0, 3)
	cp.lcs = make([]*line2CP, 0, 3)
	cp.id = id
	return cp
}

//是否包含点
func (cp *convexPolygon) isContainPoint(p point) bool {
	maxIndex := len(cp.ps) - 1
	p1 := cp.ps[maxIndex]
	p2 := cp.ps[0]
	f := (p2.X-p1.X)*(p.Y-p1.Y) - (p.X-p1.X)*(p2.Y-p1.Y)
	if f > 0 {
		return false
	}
	p1 = p2
	for i := 1; i <= maxIndex; i++ {
		p2 = cp.ps[i]
		f = (p2.X-p1.X)*(p.Y-p1.Y) - (p.X-p1.X)*(p2.Y-p1.Y)
		if f > 0 {
			return false
		}
		p1 = p2
	}

	return true
}

//是否与格子有交点
func (cp *convexPolygon) isIntersectGrid(gridNum int64, gsize int64, maxVNum int64) bool {
	minY := gridNum / maxVNum * gsize
	maxY := minY + gsize
	minX := ((gridNum - 1) % maxVNum) * gsize
	maxX := minX + gsize
	if cp.isContainPoint(point{X: minX, Y: minY}) ||
		cp.isContainPoint(point{X: maxX, Y: minY}) ||
		cp.isContainPoint(point{X: maxX, Y: maxY}) ||
		cp.isContainPoint(point{X: minX, Y: minY}) {
		return true
	}

	return false
}

//得到区域包含的格子id列表
func (cp *convexPolygon) getGrids(gsize int64, maxVNum int64) []int64 {
	maxX := cp.ps[0].X
	maxY := cp.ps[0].Y
	minX := cp.ps[0].X
	minY := cp.ps[0].Y
	t := int64(0)
	for i := 1; i < len(cp.ps); i++ {
		t = cp.ps[i].X
		if t > maxX {
			maxX = t
		}
		if t < minX {
			minX = t
		}
		t = cp.ps[i].Y
		if t > maxY {
			maxY = t
		}
		if t < minY {
			minY = t
		}
	}
	ret := make([]int64, 0, 20)
	gid := int64(0)
	p := point{}
	for x := minX / gsize * gsize; x <= maxX; x += gsize {
		for y := minY / gsize * gsize; y <= maxY; y += gsize {
			p.X, p.Y = x, y
			gid = p.getGridNum(gsize, maxVNum)
			if cp.isIntersectGrid(gid, gsize, maxVNum) {
				ret = append(ret, gid)
			}
		}
	}
	return ret
}

//构建不可穿过的线(每个区域至少有一条不可穿过的线)
func (cp *convexPolygon) makeLines() error {
	length1 := len(cp.ps)
	length2 := len(cp.lcs)
	l := line{}
	for i := 0; i < length1; i++ {
		l.sp, l.ep = cp.ps[i], cp.ps[(i+1)%length1]
		for j := 0; j < length2; j++ {
			if l.isSame(cp.lcs[j].l) {
				goto CT
			}
		}
		cp.lines = append(cp.lines, l)
	CT:
	}

	if len(cp.lines) < 1 {
		return errors.New("convex_polygon not cross line num must lager 1")
	}

	return nil
}

//相互构建区域相邻关系(区域与区域最多只有一个l2cp)
func make_l2cp(cp *convexPolygon, ocp *convexPolygon) {
	length1 := len(cp.ps)
	length2 := len(ocp.ps)
	for i := 0; i < length1; i++ {
		for j := length2 - 1; j >= 0; j-- {
			if cp.ps[i] == ocp.ps[j] {
				j1 := (j - 1 + length2) % length2
				i1 := (i + 1) % length1
				if ocp.ps[j1] == cp.ps[i1] {
					l2cp := &line2CP{
						l: line{
							sp: cp.ps[i],
							ep: cp.ps[i1],
						},
						cp: ocp,
					}
					cp.lcs = append(cp.lcs, l2cp)

					l2cp = &line2CP{
						l: line{
							sp: ocp.ps[j],
							ep: ocp.ps[j1],
						},
						cp: cp,
					}
					ocp.lcs = append(ocp.lcs, l2cp)
				}

				return
			}
		}
	}
}
