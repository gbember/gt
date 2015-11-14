// line.go
package mmap

import "math"

//线段
type line struct {
	sp point //线段起点
	ep point //线段终点
}

//得到线段穿过的格子
func (l *line) getAcossGridNums(gsize int, maxVNum int) []int {
	gnum1 := getGridNum(l.sp, gsize, maxVNum)
	gnum2 := getGridNum(l.ep, gsize, maxVNum)
	if gnum1 == gnum2 {
		return []int{gnum1}
	}
	gidList := make([]int, 0, 20)
	//在同一行
	if int(math.Abs(float64(gnum1-gnum2))) < gsize {
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
	fgsize := float32(gsize)
	gid := getGridNum(l.sp, gsize, maxVNum)
	gidList = append(gidList, gid)
	if x > 0 {
		max := float32(int(l.ep.x) / gsize * gsize)
		x = float32(int(l.sp.x)/gsize*gsize + gsize)
		//		log.Println(x, max, gsize)
		for ; x <= max; x += fgsize {
			y = tan*x + a
			gid = getGridNum(point{x: x, y: y}, gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	} else {
		min := float32(int(l.ep.x) / gsize * gsize)
		x = float32(int(l.sp.x)/gsize*gsize - gsize)
		//		log.Println(x, min, gsize)
		for ; x >= min; x -= fgsize {
			y = tan*x + a
			gid = getGridNum(point{x: x, y: y}, gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	}
	if l.ep.y-l.sp.y > 0 {
		max := float32(int(l.ep.y) / gsize * gsize)
		y = float32(int(l.sp.y)/gsize*gsize + gsize)
		//		log.Println(y, max, gsize)
		for ; y <= max; y += fgsize {
			x = (y - a) / tan
			gid = getGridNum(point{x: x, y: y}, gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	} else {
		min := float32(int(l.ep.y) / gsize * gsize)
		y = float32(int(l.sp.y)/gsize*gsize - gsize)
		//		log.Println(y, min, gsize)
		for ; y >= min; y -= fgsize {
			x = (y - a) / tan
			gid = getGridNum(point{x: x, y: y}, gsize, maxVNum)
			gidList = append(gidList, gid)
		}
	}
	return gidList
}

//是否是同一条线段(起点等于起点终点等于终点或者起点等于终点终点等于起点)
func (l *line) isEq(l1 *line) bool {
	return (l.sp == l1.sp && l.ep == l1.ep) || (l.ep == l1.sp && l.sp == l1.ep)
}

//是否穿过
func (l *line) isAcrossLine(l1 *line) bool {
	f1 := l.sp.x - l.ep.x
	f2 := l.sp.y - l.ep.y
	fC := (l1.sp.y-l.sp.y)*f1 - (l1.sp.x-l.sp.x)*f2
	fD := (l1.ep.y-l.sp.y)*f1 - (l1.ep.x-l.sp.x)*f2
	// A(x1, y1), B(x2, y2)的直线方程为：
	// f(x, y) =  (y - y1) * (x1 - x2) - (x - x1) * (y1 - y2) = 0
	if fC*fD >= 0 {
		return false
	}

	return true
}

//与N条线至少有一条交叉
func (l *line) isAcrossLines(ls []*line) bool {
	for i := 0; i < len(ls); i++ {
		if l.isAcrossLine(ls[i]) {
			return true
		}
	}
	return false
}

func (l *line) Distance() float32 {
	x := l.ep.x - l.sp.x
	y := l.ep.y - l.sp.y
	return float32(math.Sqrt(float64(x*x + y*y)))
}
