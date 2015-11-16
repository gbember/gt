// mmap.go
package mmap

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

type Map struct {
	id      uint16           //地图id
	am      map[uint32]*area //地图所有区域
	alist   []*area
	agm     map[int][]*area //格子与区域对应关系(一个格子可能在多个区域中)
	maxVNum int             //横向格子最大数量
	gsize   int             //格子大小
}

func LoadMap(bs []byte) (*Map, error) {
	pk := NewReader(bs)
	id, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	maxVNum, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	gsize, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	areaNum, err := pk.readUint16()
	if err != nil {
		return nil, err
	}
	m := new(Map)
	m.id = id
	m.maxVNum = int(maxVNum)
	m.gsize = int(gsize)
	m.am = make(map[uint32]*area)
	for i := areaNum; i > 0; i-- {
		a, err := loadArea(pk)
		if err != nil {
			return nil, err
		}
		if _, ok := m.am[a.id]; ok {
			return nil, errors.New(fmt.Sprintf("repeated area id: %d", a.id))
		}
		m.am[a.id] = a
	}
	m.init()
	return m, nil
}

//寻路
func (m *Map) FindPath(p1 point, p2 point) ([]point, bool) {
	a2 := m.getPointArea(p2)
	if a2 == nil {
		return nil, false
	}
	//	l := &line{sp: p1, ep: p2}
	//	gidList := l.getAcossGridNums(m.gsize, m.maxVNum)
	//	max := len(gidList)
	//	if max > 2 {
	//		//判断所有格子是否可以走(起点除外)
	//		isLine := true
	//		length := 0
	//		for i := max - 1; i > 0 && isLine; i-- {
	//			if aList, ok := m.agm[gidList[i]]; ok {
	//				//判断该线是否穿过这些区域不能穿过的线
	//				length = len(aList)
	//				for j := 0; j < length; j++ {
	//					if aList[j].isCrossNoPassLine(l) {
	//						isLine = false
	//						break
	//					}
	//				}
	//			} else {
	//				isLine = false
	//				break
	//			}
	//		}
	//		if isLine {
	//			log.Println("直线")
	//			return []point{p2}, true
	//		}
	//	}
	a1 := m.getPointArea(p1)
	//区域寻路
	return findPath(m, p1, a1, p2, a2)
}

//得到点所在的某个区域(可能有多个  但只返回一个 返回nil表示点不在区域中)
func (m *Map) getPointArea(p point) *area {
	gid := getGridNum(p, m.gsize, m.maxVNum)
	if as, ok := m.agm[gid]; ok {
		length := len(as)
		for i := 0; i < length; i++ {
			if as[i].isContainPoint(p) {
				return as[i]
			}
		}
	}
	return nil
}

//地图数据初始化
func (m *Map) init() {
	//1 构造格子区域关系
	m.agm = make(map[int][]*area)
	length := len(m.am)
	alist := make([]*area, 0, length)
	for _, a := range m.am {
		gidList := a.getGrids(m.gsize, m.maxVNum)
		for i := 0; i < len(gidList); i++ {
			aList, ok := m.agm[gidList[i]]
			if ok {
				aList = append(aList, a)
			} else {
				aList = make([]*area, 0, 10)
				aList = append(aList, a)
			}
			m.agm[gidList[i]] = aList
		}
		alist = append(alist, a)
	}
	//2 构造区域不可穿过线与构造区域与区域关系
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			alist[i].makeLineAndRela(alist[j])
		}
	}

	m.alist = alist
	m.am = nil
}

func Test() {
	m := new(Map)
	m.gsize = 50
	m.id = 1
	m.maxVNum = 1000
	m.am = make(map[uint32]*area)
	max1 := uint32(25)
	max2 := uint32(25)
	for k := uint32(1); k <= max1; k++ {
		for i := uint32(1); i <= max2; i++ {
			a := new(area)
			a.id = k*max1 + i
			a.points = make([]point, 0, 4)
			a.allLines = make([]*line, 0, 4)
			a.lineMap = make(map[*line]bool)
			a.areaMap = make(map[*line]*area)
			p := point{x: float32(i) * 8, y: float32(k) * 8}
			a.points = append(a.points, p)

			p = point{x: float32(i+1) * 8, y: float32(k) * 8}
			a.points = append(a.points, p)
			l := &line{sp: a.points[0], ep: p}
			a.allLines = append(a.allLines, l)

			p = point{x: float32(i+1) * 8, y: float32(k+1) * 8}
			a.points = append(a.points, p)
			l = &line{sp: a.points[1], ep: p}
			a.allLines = append(a.allLines, l)

			p = point{x: float32(i) * 8, y: float32(k+1) * 8}
			a.points = append(a.points, p)
			l = &line{sp: a.points[2], ep: p}
			a.allLines = append(a.allLines, l)

			l = &line{sp: a.points[3], ep: a.points[0]}
			a.allLines = append(a.allLines, l)
			m.am[a.id] = a
		}
	}
	m.init()

	fn := "tt.cpuprof"
	f, err := os.Create(fn)
	isProf := false
	if err == nil {
		err = pprof.StartCPUProfile(f)
		if err == nil {
			isProf = true
		}
	}

	points, ok := m.FindPath(point{x: 28, y: 28}, point{x: 400, y: 305})
	if ok {
		log.Println(points)
		startTime := time.Now()
		var max int64 = 1000000
		for i := max; i > 0; i-- {
			m.FindPath(point{x: 28, y: 28}, point{x: 400, y: 305})
		}
		td := time.Since(startTime)
		log.Println(td)
		log.Println(td.Nanoseconds() / max)
	} else {
		log.Println("终点不可达")
	}

	if isProf {
		pprof.StopCPUProfile()
	}
}
