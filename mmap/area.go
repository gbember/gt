// area.go
package mmap

import "errors"

//区域
type area struct {
	id       uint32          //区域ID
	points   []point         //区域定点
	allLines []*line         //区域所有线段
	lineMap  map[*line]bool  //区域不可穿过线段(为nil)
	lines    []*line         //区域不可穿过线段
	areaMap  map[*line]*area //区域相邻区域(线相连)
}

func loadArea(pk *Packet) (*area, error) {
	id, err := pk.readUint32()
	if err != nil {
		return nil, err
	}
	pointNum, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	if pointNum < 3 {
		return nil, errors.New("area point num must lager 3")
	}
	a := new(area)
	a.id = id
	a.points = make([]point, 0, pointNum)
	a.allLines = make([]*line, 0, pointNum)
	a.lineMap = make(map[*line]bool)
	a.areaMap = make(map[*line]*area)
	//第一个
	p := point{}
	p.x, err = pk.readInt32()
	if err != nil {
		return nil, err
	}
	p.y, err = pk.readInt32()
	if err != nil {
		return nil, err
	}
	a.points = append(a.points, p)
	for i := uint16(1); i < pointNum; i++ {
		p = point{}
		p.x, err = pk.readInt32()
		if err != nil {
			return nil, err
		}
		p.y, err = pk.readInt32()
		if err != nil {
			return nil, err
		}
		a.points = append(a.points, p)
		l := &line{sp: a.points[i-1], ep: p}
		a.lineMap[l] = true
		a.allLines = append(a.allLines, l)
	}
	l := &line{sp: a.points[pointNum-1], ep: a.points[0]}
	a.lineMap[l] = true
	a.allLines = append(a.allLines, l)
	return a, nil
}

//构造区域不可穿过线与构造区域与区域关系
func (a *area) makeLineAndRela(a1 *area) {
	for i := 0; i < len(a.allLines); i++ {
		for j := 0; j < len(a1.allLines); j++ {
			if a.allLines[i].isEq(a1.allLines[j]) {
				delete(a.lineMap, a.allLines[i])
				delete(a1.lineMap, a1.allLines[j])
				a.areaMap[a.allLines[i]] = a1
				a1.areaMap[a1.allLines[j]] = a
				break
			}
		}
	}
}

//是否穿过不能穿过的线
func (a *area) isCrossNoPassLine(l *line) bool {
	for l1, _ := range a.lineMap {
		if l1.isAcrossLine(l) {
			return true
		}
	}
	return false
}

//得到区域格子
func (a *area) getGrids(gsize int32, maxVNum int32) []int32 {
	var maxX int32 = a.points[0].x
	var maxY int32 = a.points[0].y
	var minX int32 = a.points[0].x
	var minY int32 = a.points[0].y
	var t int32 = 0
	for i := 1; i < len(a.points); i++ {
		t = a.points[i].x
		if t > maxX {
			maxX = t
		}
		if t < minX {
			minX = t
		}
		t = a.points[i].y
		if t > maxY {
			maxY = t
		}
		if t < minY {
			minY = t
		}
	}
	ret := make([]int32, 0, 20)
	var gid int32 = 0
	for x := minX / gsize * gsize; x <= maxX; x += gsize {
		for y := minY / gsize * gsize; y <= maxY; y += gsize {
			gid = getGridNum(point{x: x, y: y}, gsize, maxVNum)
			if a.isContainGrid(gid, gsize, maxVNum) {
				ret = append(ret, gid)
			}
		}
	}
	return ret
}

//是否包含格子(4条线至少有一条穿过区域)
func (a *area) isContainGrid(gridNum int32, gsize int32, maxVNum int32) bool {
	minY := gridNum / maxVNum * gsize
	maxY := minY + gsize
	minX := ((gridNum - 1) % maxVNum) * gsize
	maxX := minX + gsize

	l := &line{sp: point{x: minX, y: minY}, ep: point{x: maxX, y: minY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	l = &line{sp: point{x: maxX, y: minY}, ep: point{x: maxX, y: maxY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	l = &line{sp: point{x: maxX, y: maxY}, ep: point{x: minX, y: maxY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	l = &line{sp: point{x: minX, y: maxY}, ep: point{x: minX, y: minY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	return false
}

//区域是否包含某个点
func (a *area) isContainPoint(p point) bool {
	//	l := &line{p, point{math.MaxInt32, math.MaxInt32}}
	//	return l.isAcrossLines(a.allLines)

	maxIndex := len(a.points) - 1
	p1 := a.points[maxIndex]
	p2 := a.points[0]

	f := (p2.x-p1.x)*(p.y-p1.y) - (p.x-p1.x)*(p2.y-p1.y)
	if f > 0 {
		return false
	}
	p1 = a.points[0]
	for i := 1; i <= maxIndex; i++ {
		p2 = a.points[i]
		f = (p2.x-p1.x)*(p.y-p1.y) - (p.x-p1.x)*(p2.y-p1.y)
		if f > 0 {
			return false
		}
		p1 = p2
	}

	return true
}

//func mul(p point, p1 point, p2 point) int32 {
//	return (p2.x-p1.x)*(p.y-p1.y) - (p.x-p1.x)*(p2.y-p1.y)
//}
