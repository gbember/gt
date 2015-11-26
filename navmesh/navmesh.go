// navmesh.go
package navmesh

import(
	"math"
)

type NavMesh struct{
	//格子大小(用于定位点在某个多边形区域)
	gsize int64
	//地图宽度
	width int64
	//地图高度
	heigth int64
	//格子大小和地图宽度计算出来的横向最多格子数
	maxVNum int64
}

//得到点所在的格子编号
func (nm *NavMesh)getGridNum(p point)int64{
	return p.x/nm.gsize + 1 + p.y/nm.gsize*nm.maxVNum
}

//得到线段穿过的格子id列表
func (nm *NavMesh) getAcossGridNums(l line) []int64 {
	gnum1 := nm.getGridNum(l.sp)
	gnum2 := nm.getGridNum(l.ep)
	if gnum1 == gnum2 {
		return []int64{gnum1}
	}
	gidList := make([]int64, 0, 20)
	gsize := nm.gsize
	maxVNum := nm.maxVNum
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
	x := l.ep.x - l.sp.x
	y := l.ep.y - l.sp.y
	tan := y / x
	a := l.ep.y - tan*l.ep.x
	gid := nm.getGridNum(l.sp)
	gidList = append(gidList, gid)
	if x > 0 {
		max := l.ep.x / gsize * gsize
		x = l.sp.x/gsize*gsize + gsize
		for ; x <= max; x += gsize {
			y = tan*x + a
			gid = nm.getGridNum(point{x: x, y: y})
			gidList = append(gidList, gid)
		}
	} else {
		min := l.ep.x / gsize * gsize
		x = l.sp.x/gsize*gsize - gsize
		for ; x >= min; x -= gsize {
			y = tan*x + a
			gid = nm.getGridNum(point{x: x, y: y})
			gidList = append(gidList, gid)
		}
	}
	if l.ep.y-l.sp.y > 0 {
		max := l.ep.y / gsize * gsize
		y = l.sp.y/gsize*gsize + gsize
		for ; y <= max; y += gsize {
			x = (y - a) / tan
			gid = nm.getGridNum(point{x: x, y: y})
			gidList = append(gidList, gid)
		}
	} else {
		min := l.ep.y / gsize * gsize
		y = l.sp.y/gsize*gsize - gsize
		for ; y >= min; y -= gsize {
			x = (y - a) / tan
			gid = nm.getGridNum(point{x: x, y: y})
			gidList = append(gidList, gid)
		}
	}
	return gidList
}

//得到区域包含的格子id列表
func (nm *NavMesh) getGrids(cp *ConvexPolygon) []int64 {
	maxX := cp.ps[0].x
	maxY := cp.ps[0].y
	minX := cp.ps[0].x
	minY := cp.ps[0].y
	gsize:=nm.gsize
	t := int64(0)
	for i := 1; i < len(cp.ps); i++ {
		t = cp.ps[i].x
		if t > maxX {
			maxX = t
		}
		if t < minX {
			minX = t
		}
		t = cp.ps[i].y
		if t > maxY {
			maxY = t
		}
		if t < minY {
			minY = t
		}
	}
	ret := make([]int64, 0, 20)
	gid := int64(0)
	for x := minX / gsize * gsize; x <= maxX; x += gsize {
		for y := minY / gsize * gsize; y <= maxY; y += gsize {
			gid = nm.getGridNum(point{x: x, y: y})
			if cp.isContainGrid(gid) {
				ret = append(ret, gid)
			}
		}
	}
	return ret
}
