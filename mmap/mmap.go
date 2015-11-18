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
	if a1 == a2 {
		return []point{p2}, true
	}

	if fppList, ok := a1.fpMap[a2]; ok {
		length := len(fppList)
		if length == 0 {
			return nil, false
		}
		l := &line{}
		size := float64(0)
		minSize := float64(9999999999999999999999999999)
		minIndex := 0
		for i := 0; i < length; i++ {
			l.sp = p1
			l.ep = fppList[i].ps[0]
			size = fppList[i].size + l.Distance()
			l.sp = p2
			l.ep = fppList[i].ps[len(fppList[i].ps)-1]
			size += l.Distance()
			if size <= minSize {
				minSize = size
				minIndex = i
			}
		}
		lps := len(fppList[minIndex].ps)
		if p2 == fppList[minIndex].ps[lps-1] {
			ps := make([]point, lps, lps)
			copy(ps, fppList[minIndex].ps)
			return ps, true
		} else {
			ps := make([]point, lps+1, lps+1)
			copy(ps, fppList[minIndex].ps)
			ps[lps] = p2
			return ps, true
		}

	}
	return nil, false
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

	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if i != j {

				//				if alist[i].id == 1 && alist[j].id == 30 {
				//					log.Println("11111111111")
				//				}

				//				if alist[i].id == 30 && alist[j].id == 1 {
				//					log.Println("2222222222222")
				//				}

				if _, ok := alist[i].fpMap[alist[j]]; !ok {
					length := len(alist[i].areaMap) * 2 * len(alist[j].areaMap) * 2
					pps1 := make([]*FPointPath, 0, length)
					pps2 := make([]*FPointPath, 0, length)

					for l1, a1 := range alist[i].areaMap {
						for l2, _ := range alist[j].areaMap {
							fpp := findPath(m, l1.sp, a1, l2.sp, alist[j])
							if fpp != nil {
								pps1 = append(pps1, fpp)
								pps2 = append(pps2, fpp.reverse())
							}
							fpp = findPath(m, l1.sp, a1, l2.ep, alist[j])
							if fpp != nil {
								pps1 = append(pps1, fpp)
								pps2 = append(pps2, fpp.reverse())
							}
							fpp = findPath(m, l1.ep, a1, l2.sp, alist[j])
							if fpp != nil {
								pps1 = append(pps1, fpp)
								pps2 = append(pps2, fpp.reverse())
							}
							fpp = findPath(m, l1.ep, a1, l2.ep, alist[j])
							if fpp != nil {
								pps1 = append(pps1, fpp)
								pps2 = append(pps2, fpp.reverse())
							}
						}
					}
					alist[i].fpMap[alist[j]] = pps1
					alist[j].fpMap[alist[i]] = pps2
				}
			}
		}
	}
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
	destP := point{170, 889}
	orgiP := point{x: 101, y: 86}

	ps, ok := m.FindPath(orgiP, destP)
	if ok {
		log.Println(ps)
		startTime := time.Now()
		var max int64 = 1000000
		for i := max; i > 0; i-- {
			m.FindPath(orgiP, destP)
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
	a.fpMap = make(map[*area][]*FPointPath)

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
	addArea(m.am, 12, point{236, 371}, point{213, 350}, point{140, 380}, point{170, 397})
	addArea(m.am, 13, point{170, 397}, point{140, 380}, point{114, 424}, point{164, 435})
	addArea(m.am, 14, point{164, 435}, point{114, 424}, point{68, 458}, point{119, 463})
	addArea(m.am, 15, point{119, 463}, point{68, 458}, point{123, 515}, point{164, 513})
	addArea(m.am, 16, point{164, 513}, point{123, 515}, point{122, 552}, point{216, 549})
	addArea(m.am, 17, point{164, 513}, point{216, 549}, point{268, 503}, point{238, 491})
	addArea(m.am, 18, point{122, 552}, point{123, 515}, point{75, 507}, point{71, 524}, point{74, 544})
	addArea(m.am, 19, point{74, 544}, point{71, 524}, point{41, 540}, point{18, 579}, point{48, 576})
	addArea(m.am, 20, point{48, 576}, point{18, 579}, point{46, 619}, point{91, 619}, point{86, 599})
	addArea(m.am, 21, point{86, 599}, point{91, 619}, point{148, 601}, point{136, 585})
	addArea(m.am, 22, point{216, 549}, point{122, 552}, point{136, 585}, point{148, 601}, point{203, 585})
	addArea(m.am, 23, point{216, 549}, point{203, 585}, point{213, 606}, point{246, 587})
	addArea(m.am, 24, point{246, 587}, point{213, 606}, point{259, 627}, point{280, 611})
	addArea(m.am, 25, point{280, 611}, point{259, 627}, point{302, 709}, point{334, 694})
	addArea(m.am, 26, point{334, 694}, point{302, 709}, point{340, 758}, point{394, 743})
	addArea(m.am, 27, point{394, 743}, point{340, 758}, point{404, 835}, point{456, 836}, point{447, 795})
	addArea(m.am, 28, point{404, 835}, point{340, 758}, point{306, 814}, point{335, 898})
	addArea(m.am, 29, point{335, 898}, point{306, 814}, point{226, 874}, point{223, 899})
	addArea(m.am, 30, point{223, 899}, point{226, 874}, point{123, 874}, point{118, 896})

	m.init()
	return m
}
