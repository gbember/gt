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
	agm     map[int32][]*area //格子与区域对应关系(一个格子可能在多个区域中)
	maxVNum int32             //横向格子最大数量
	gsize   int32             //格子大小
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
	m.maxVNum = int32(maxVNum)
	m.gsize = int32(gsize)
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
	m.agm = make(map[int32][]*area)
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
	m := loadMap_1()

	fn := "tt.cpuprof"
	f, err := os.Create(fn)
	isProf := false
	if err == nil {
		err = pprof.StartCPUProfile(f)
		if err == nil {
			isProf = true
		}
	}

	points, ok := m.FindPath(point{x: 101, y: 86}, point{x: 328, y: 324})
	if ok {
		log.Println(points)
		startTime := time.Now()
		var max int64 = 10000
		for i := max; i > 0; i-- {
			m.FindPath(point{x: 101, y: 86}, point{x: 328, y: 324})
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

func addArea(am map[uint32]*area, id uint32, ps ...point) {
	a := new(area)
	a.id = id
	length := len(ps)
	if length < 3 {
		panic("area point num must lager 3")
	}
	a.points = ps
	a.allLines = make([]*line, 0, length)
	a.lineMap = make(map[*line]bool)
	a.areaMap = make(map[*line]*area)

	p := ps[0]
	l := &line{ps[length-1], p}
	a.allLines = append(a.allLines, l)

	for i := 1; i < length; i++ {
		l = &line{p, ps[i]}
		a.allLines = append(a.allLines, l)
		p = ps[i]
	}

	am[id] = a

}

func loadMap_1() *Map {
	m := new(Map)
	m.gsize = 5
	m.id = 1
	m.maxVNum = 1000
	m.am = make(map[uint32]*area)

	addArea(m.am, 1, point{75, 95}, point{118, 97}, point{188, 50}, point{190, 23})
	addArea(m.am, 2, point{190, 23}, point{188, 50}, point{294, 59}, point{307, 39}, point{207, 22})
	addArea(m.am, 3, point{270, 22}, point{307, 39}, point{318, 2}, point{282, 3})
	addArea(m.am, 4, point{307, 39}, point{294, 59}, point{354, 101}, point{381, 65})
	addArea(m.am, 5, point{381, 65}, point{354, 101}, point{408, 97}, point{364, 54}, point{566, 2}, point{476, 3})
	addArea(m.am, 6, point{408, 97}, point{354, 101}, point{380, 134}, point{425, 125})
	addArea(m.am, 7, point{425, 125}, point{380, 134}, point{368, 182}, point{420, 186})
	addArea(m.am, 8, point{420, 186}, point{368, 182}, point{349, 216}, point{403, 217})
	addArea(m.am, 9, point{403, 217}, point{349, 216}, point{291, 280}, point{347, 291})
	addArea(m.am, 10, point{347, 291}, point{291, 280}, point{246, 318}, point{347, 340})
	addArea(m.am, 11, point{347, 340}, point{246, 318}, point{213, 350}, point{236, 371}, point{358, 382})

	m.init()
	return m
}
