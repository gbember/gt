// mmap.go
package mmap

import (
	"errors"
	"fmt"
	"log"
)

type Map struct {
	id      uint16           //地图id
	am      map[uint16]*area //地图所有区域
	agm     map[int][]uint16 //格子与区域编号对应关系(一个格子可能在多个区域中)
	maxVNum int              //横向格子最大数量
	gsize   int              //格子大小
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
	m.am = make(map[uint16]*area)
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
func (m *Map) FindPath(p1 point, p2 point) []point {
	egid := getGridNum(p2, m.gsize, m.maxVNum)
	//终点可以走
	if _, ok := m.agm[egid]; ok {
		l := &line{sp: p1, ep: p2}
		gidList := l.getAcossGridNums(m.gsize, m.maxVNum)
		max := len(gidList)
		if max > 2 {
			//判断所有格子是否可以走(起点除外)
			isLine := true
			for i := 1; i < max && isLine; i++ {
				am := make(map[uint16]bool)
				if aidList, ok := m.agm[gidList[i]]; ok {
					//判断该线是否穿过这些区域不能穿过的线
					aid := uint16(0)
					for j := 0; j < len(aidList); j++ {
						aid = aidList[j]
						if !am[aid] {
							if m.am[aid].isCrossNoPassLine(l) {
								isLine = false
								break
							}
							am[aid] = true
						}
					}
				} else {
					isLine = false
					break
				}
			}
			if isLine {
				return []point{p2}
			}
			//TODO 区域寻路
			return nil
		} else {
			return []point{p2}
		}
	}
	return nil
}

//地图数据初始化
func (m *Map) init() {
	//1 构造格子区域关系
	m.agm = make(map[int][]uint16)
	length := len(m.am)
	alist := make([]*area, 0, length)
	for _, a := range m.am {
		aid := a.id
		gidList := a.getGrids(m.gsize, m.maxVNum)
		for i := 0; i < len(gidList); i++ {
			aidList, ok := m.agm[gidList[i]]
			if ok {
				aidList = append(aidList, aid)
			} else {
				aidList = make([]uint16, 0, 10)
				aidList = append(aidList, aid)
			}
			m.agm[gidList[i]] = aidList
		}
		alist = append(alist, a)
	}
	//2 构造区域不可穿过线与构造区域与区域关系
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			alist[i].makeLineAndRela(alist[j])
		}
	}
}

func Test() {
	m := new(Map)
	m.gsize = 5
	m.id = 1
	m.maxVNum = 1000
	m.am = make(map[uint16]*area)
	for k := uint16(1); k <= 10; k++ {
		for i := uint16(1); i <= 10; i++ {
			a := new(area)
			a.id = i * k
			a.points = make([]point, 0, 4)
			a.allLines = make([]*line, 0, 4)
			a.lines = make(map[*line]bool)
			a.areaMap = make(map[*area]bool)
			p := point{x: float32(i) * 28, y: float32(k) * 28}
			a.points = append(a.points, p)

			p = point{x: float32(i+1) * 28, y: float32(k) * 28}
			a.points = append(a.points, p)
			l := &line{sp: a.points[0], ep: p}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)

			p = point{x: float32(i+1) * 28, y: float32(k+1) * 28}
			a.points = append(a.points, p)
			l = &line{sp: a.points[1], ep: p}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)

			p = point{x: float32(i) * 28, y: float32(k+1) * 28}
			a.points = append(a.points, p)
			l = &line{sp: a.points[2], ep: p}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)

			l = &line{sp: a.points[4-1], ep: a.points[0]}
			a.lines[l] = true
			a.allLines = append(a.allLines, l)
			m.am[a.id] = a
		}
	}
	m.init()
	points := m.FindPath(point{x: 30, y: 30}, point{x: 250, y: 270})
	log.Println(points)
}
