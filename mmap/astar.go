// astar.go
package mmap

import "container/heap"

type APoint struct {
	p        point
	a        *area
	size     float32
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
func (ol *OpenList) getMin() float32 {
	var min float32 = 9999999999
	for i := 0; i < ol.length; i++ {
		if ol.apList[i].size < min {
			min = ol.apList[i].size
		}
	}
	return min
}

func findPath(m *Map, p1 point, a1 *area, p2 point, a2 *area) ([]point, bool) {
	mas := &MapAStar{
		openList: new(OpenList),
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
	for mas.openList.length > 0 {
		apx := heap.Pop(mas.openList)
		ap = apx.(*APoint)
		if ap.p == mas.destP {
			ps := make([]point, 0, ap.length)
			for ap != nil {
				ps = append(ps, ap.p)
				ap = ap.parentAP
			}
			//TODO 反序
			return ps, true
		}
		mas.addCloseList(ap.p)
		mas.addNextAPOpenList(ap)

	}
	return nil, false
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
		if ap.p != l.sp && !mas.closeMap[l.sp] {
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
		if ap.p != l.ep && !mas.closeMap[l.ep] {
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

func (mas *MapAStar) addCloseList(p point) {
	mas.closeMap[p] = true
}
