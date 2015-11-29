// convex_polygon.go
package navmesh

import "errors"

//凸多边形
type convexPolygon struct {
	id     int        //区域编号
	pindex []uint16   //多边形的所有顶点(逆时针排序)
	lines  []line     //该多边形不能穿过的线
	lcs    []*line2CP //相邻区域
}

//多边形相邻关系
type line2CP struct {
	spindex uint16
	epindex uint16
	cp      *convexPolygon //区域
}

func newConvexPolygon(id int) *convexPolygon {
	cp := new(convexPolygon)
	cp.lines = make([]line, 0, 3)
	cp.lcs = make([]*line2CP, 0, 3)
	cp.id = id
	return cp
}

//是否包含点
func (cp *convexPolygon) isContainPoint(nm *NavMesh, p Point) bool {
	length := len(cp.pindex)

	p1 := Point{}
	p2 := Point{}
	for i := 0; i < length; i++ {
		p1 = nm.points[cp.pindex[i]]
		p2 = nm.points[cp.pindex[(i+1)%length]]
		if (p2.X-p1.X)*(p.Y-p1.Y)-(p.X-p1.X)*(p2.Y-p1.Y) > 0 {
			return false
		}
	}
	return true
}

//是否与格子有交点
func (cp *convexPolygon) isIntersectGrid(nm *NavMesh, gridNum int64) bool {
	minY := gridNum / nm.maxVNum * nm.gsize
	maxY := minY + nm.gsize
	minX := ((gridNum - 1) % nm.maxVNum) * nm.gsize
	maxX := minX + nm.gsize
	if cp.isContainPoint(nm, Point{X: minX, Y: minY}) ||
		cp.isContainPoint(nm, Point{X: maxX, Y: minY}) ||
		cp.isContainPoint(nm, Point{X: maxX, Y: maxY}) ||
		cp.isContainPoint(nm, Point{X: minX, Y: minY}) {
		return true
	}

	return false
}

//得到区域包含的格子id列表
func (cp *convexPolygon) getGrids(nm *NavMesh) []int64 {
	index := cp.pindex[0]
	maxX := nm.points[index].X
	maxY := nm.points[index].Y
	minX := nm.points[index].X
	minY := nm.points[index].Y
	t := int64(0)
	for i := 1; i < len(cp.pindex); i++ {
		index = cp.pindex[i]
		t = nm.points[index].X
		if t > maxX {
			maxX = t
		}
		if t < minX {
			minX = t
		}
		t = nm.points[index].Y
		if t > maxY {
			maxY = t
		}
		if t < minY {
			minY = t
		}
	}
	ret := make([]int64, 0, 20)
	gid := int64(0)
	p := Point{}
	gsize := nm.gsize
	maxVNum := nm.maxVNum
	for x := minX / gsize * gsize; x <= maxX; x += gsize {
		for y := minY / gsize * gsize; y <= maxY; y += gsize {
			p.X, p.Y = x, y
			gid = p.getGridNum(gsize, maxVNum)
			if cp.isIntersectGrid(nm, gid) {
				ret = append(ret, gid)
			}
		}
	}
	return ret
}

//构建不可穿过的线(每个区域至少有一条不可穿过的线)
func (cp *convexPolygon) makeLines(nm *NavMesh) error {
	length1 := len(cp.pindex)
	length2 := len(cp.lcs)
	l := line{}
	l1 := line{}
	for i := 0; i < length1; i++ {
		l.sp, l.ep = nm.points[cp.pindex[i]], nm.points[cp.pindex[(i+1)%length1]]
		for j := 0; j < length2; j++ {
			l1.sp, l1.ep = nm.points[cp.lcs[j].spindex], nm.points[cp.lcs[j].epindex]
			if l.isSame(l1) {
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
func (nm *NavMesh) make_l2cp(cp *convexPolygon, ocp *convexPolygon) {
	length1 := len(cp.pindex)
	length2 := len(ocp.pindex)
	for i := 0; i < length1; i++ {
		for j := length2 - 1; j >= 0; j-- {
			if cp.pindex[i] == ocp.pindex[j] {
				j1 := (j - 1 + length2) % length2
				i1 := (i + 1) % length1
				if ocp.pindex[j1] == cp.pindex[i1] {
					l2cp := &line2CP{
						spindex: cp.pindex[i],
						epindex: cp.pindex[i1],
						cp:      ocp,
					}
					cp.lcs = append(cp.lcs, l2cp)

					l2cp = &line2CP{
						spindex: ocp.pindex[j],
						epindex: ocp.pindex[j1],
						cp:      cp,
					}
					ocp.lcs = append(ocp.lcs, l2cp)
				}

				return
			}
		}
	}
}
