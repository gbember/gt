// astar.go
package mmap

import "container/heap"

type APoint struct {
	p        point
	a        *area
	size     float64
	length   int
	parentAP *APoint
}

type MapAStar struct {
	//开发列表
	openList *OpenList
	//关闭列表
	closeMap map[point]bool
	//起点
	origP point
	//起点区域
	origA *area
	//终点
	destP point
	//终点区域
	destA *area
}

type OpenList struct {
	length int
	apList []*APoint
}

func (ol *OpenList) Len() int           { return ol.length }
func (ol *OpenList) Less(i, j int) bool { return ol.apList[i].size < ol.apList[j].size }
func (ol *OpenList) Swap(i, j int)      { ol.apList[i], ol.apList[j] = ol.apList[j], ol.apList[i] }
func (ol *OpenList) Push(x interface{}) {
	ol.length++
	ol.apList = append(ol.apList, x.(*APoint))
}
func (ol *OpenList) Pop() interface{} {
	x := ol.apList[ol.length-1]
	ol.length--
	ol.apList = ol.apList[:ol.length]
	return x
}
func (ol *OpenList) getMin() float64 {
	var min float64 = 99999999999999
	for i := 0; i < ol.length; i++ {
		if ol.apList[i].size < min {
			min = ol.apList[i].size
		}
	}
	return min
}

type FPointPath struct {
	size float64
	ps   []point
}

//反转路径点
func (fpp *FPointPath) reverse() *FPointPath {
	nfpp := new(FPointPath)
	nfpp.size = fpp.size
	l := len(fpp.ps)
	nfpp.ps = make([]point, l, l)
	maxIndex := l - 1
	for i := 0; i < l; i++ {
		nfpp.ps[i] = fpp.ps[maxIndex-i]
	}
	return nfpp
}

func findPath(m *Map, p1 point, a1 *area, p2 point, a2 *area) *FPointPath {
	mas := &MapAStar{
		openList: &OpenList{
			apList: make([]*APoint, 0, 20),
		},
		closeMap: make(map[point]bool),
		origP:    p1,
		origA:    a1,
		destP:    p2,
		destA:    a2,
	}
	heap.Init(mas.openList)
	ap := &APoint{
		a:      a1,
		p:      p1,
		length: 1,
	}
	heap.Push(mas.openList, ap)
	var apx interface{}
	for mas.openList.length > 0 {
		apx = heap.Pop(mas.openList)
		ap = apx.(*APoint)
		if ap.p == mas.destP {
			fpp := new(FPointPath)
			ps := make([]point, 0, ap.length)
			for ap != nil {
				fpp.size += ap.size
				ps = append(ps, ap.p)
				ap = ap.parentAP
			}
			fpp.ps = ps
			return fpp.reverse()
		}
		mas.addNextAPOpenList(ap)
		mas.addCloseList(ap.p)
	}
	return nil
}

func (mas *MapAStar) addNextAPOpenList(ap *APoint) {
	var ap1 *APoint
	var li *line
	if ap.a.id == mas.destA.id {
		li = &line{ap.p, mas.destP}
		ap1 = &APoint{
			a:        ap.a,
			p:        mas.destP,
			size:     ap.size + li.Distance(),
			parentAP: ap,
			length:   ap.length + 1,
		}
		heap.Push(mas.openList, ap1)
		return
	}

	for l, a := range ap.a.areaMap {
		if !mas.closeMap[l.sp] {
			if l.sp == ap.p {
				ap.a = a
				heap.Push(mas.openList, ap)
			} else {
				li = &line{l.sp, ap.p}
				ap1 = &APoint{
					a:        a,
					p:        l.sp,
					size:     ap.size + li.Distance(),
					parentAP: ap,
					length:   ap.length + 1,
				}
				heap.Push(mas.openList, ap1)
			}
		}
		if !mas.closeMap[l.ep] {
			if l.ep == ap.p {
				ap.a = a
				heap.Push(mas.openList, ap)
			} else {
				li = &line{l.ep, ap.p}
				ap1 = &APoint{
					a:        a,
					p:        l.ep,
					size:     ap.size + li.Distance(),
					parentAP: ap,
					length:   ap.length + 1,
				}
				heap.Push(mas.openList, ap1)
			}
		}
	}

}

func (mas *MapAStar) addCloseList(p point) {
	mas.closeMap[p] = true
}
