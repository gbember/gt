// area.go
package mmap

import (
	"errors"
)

//区域
type area struct {
	id       uint16         //区域ID
	points   []point        //区域定点
	allLines []*line        //区域所有线段
	lines    map[*line]bool //区域不可穿过线段
	areaMap  map[*area]bool //区域相邻区域
}

func loadArea(pk *Packet) (*area, error) {
	id, err := pk.readUint16()
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
	a.lines = make(map[*line]bool)
	a.areaMap = make(map[*area]bool)
	//第一个
	p := point{}
	p.x, err = pk.readFloat32()
	if err != nil {
		return nil, err
	}
	p.y, err = pk.readFloat32()
	if err != nil {
		return nil, err
	}
	a.points = append(a.points, p)
	for i := uint16(1); i < pointNum; i++ {
		p = point{}
		p.x, err = pk.readFloat32()
		if err != nil {
			return nil, err
		}
		p.y, err = pk.readFloat32()
		if err != nil {
			return nil, err
		}
		a.points = append(a.points, p)
		l := &line{sp: a.points[i-1], ep: p}
		a.lines[l] = true
		a.allLines = append(a.allLines, l)
	}
	l := &line{sp: a.points[pointNum-1], ep: a.points[0]}
	a.lines[l] = true
	a.allLines = append(a.allLines, l)
	return a, nil
}

//构造区域不可穿过线与构造区域与区域关系
func (a *area) makeLineAndRela(a1 *area) {
	for i := 0; i < len(a.allLines); i++ {
		for j := 0; j < len(a1.allLines); j++ {
			if a.allLines[i].isEq(a1.allLines[j]) {
				delete(a.lines, a.allLines[i])
				delete(a1.lines, a1.allLines[j])
				a.areaMap[a1] = true
				a1.areaMap[a] = true
				break
			}
		}
	}
}

//是否穿过不能穿过的线
func (a *area) isCrossNoPassLine(l *line) bool {
	for l1, _ := range a.lines {
		if l1.isAcrossLine(l) {
			return true
		}
	}
	return false
}

//得到区域格子
func (a *area) getGrids(gsize int, maxVNum int) []int {
	var maxX float32 = a.points[0].x
	var maxY float32 = a.points[0].y
	var minX float32 = a.points[0].x
	var minY float32 = a.points[0].y
	var t float32 = 0
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
	ret := make([]int, 0, 20)
	fgsize := float32(gsize)
	var gid int = 0
	for x := minX; x <= maxX; x += fgsize {
		for y := minY; y < maxY; y += fgsize {
			gid = getGridNum(point{x: x, y: y}, gsize, maxVNum)
			if a.isContainGrid(gid, gsize, maxVNum) {
				ret = append(ret, gid)
			}
		}
	}
	return ret
}

//是否包含格子(4条线至少有一条穿过区域)
func (a *area) isContainGrid(gridNum int, gsize int, maxVNum int) bool {
	minY := gridNum / maxVNum * gsize
	maxY := minY + gsize
	minX := ((gridNum - 1) % maxVNum) * gsize
	maxX := minX + gsize

	fminY := float32(minY)
	fmaxY := float32(maxY)
	fminX := float32(minX)
	fmaxX := float32(maxX)
	l := &line{sp: point{x: fminX, y: fminY}, ep: point{x: fmaxX, y: fminY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	l = &line{sp: point{x: fmaxX, y: fminY}, ep: point{x: fmaxX, y: fmaxY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	l = &line{sp: point{x: fmaxX, y: fmaxY}, ep: point{x: fminX, y: fmaxY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	l = &line{sp: point{x: fminX, y: fmaxY}, ep: point{x: fminX, y: fminY}}
	if l.isAcrossLines(a.allLines) {
		return true
	}
	return true
}
