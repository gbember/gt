// line.go
package mmap

//线段
type line struct {
	sp point //线段起点
	ep point //线段终点
}

//得到线段穿过的格子
func (l *line) getAcossGridNums(gsize int, maxVNum int) []int {
	return getGridNums(l, gsize, maxVNum)
}

//是否交叉
func (l *line) isCross(l1 *line) bool {
	f1 := l.sp.x - l.ep.x
	f2 := l.sp.y - l.ep.y
	fC := (l1.sp.y-l.sp.y)*f1 - (l1.sp.x-l.sp.x)*f2
	fD := (l1.ep.y-l.sp.y)*f1 - (l1.ep.x-l.sp.x)*f2
	// A(x1, y1), B(x2, y2)的直线方程为：
	// f(x, y) =  (y - y1) * (x1 - x2) - (x - x1) * (y1 - y2) = 0
	if fC*fD > 0 {
		return false
	}

	return true
}
