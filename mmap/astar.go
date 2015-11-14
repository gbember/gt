// astar.go
package mmap

type APoint struct {
	p        point
	a        *area
	size     float32
	length   int
	parentAP *APoint
}

type MapAStar struct {
	//开发列表
	openList []*APoint
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

func findPath(m *Map, p1 point, a1 *area, p2 point, a2 *area) []point {
	mas := &MapAStar{
		openList: make([]*APoint, 0, 20),
		closeMap: make(map[point]bool),
		origP:    p1,
		origA:    a1,
		destP:    p2,
		destA:    a2,
	}
	ap := &APoint{
		a:      a1,
		p:      p1,
		length: 1,
	}
	mas.openList = append(mas.openList, ap)
	for {
		ap = mas.getMinSizeAPoint()
		if ap == nil {
			//没有路径可达
			return nil
		}
		if ap.a.id == mas.destA.id {
			ps := make([]point, 0, ap.length+1)
			for ap != nil {
				ps = append(ps, ap.p)
				ap = ap.parentAP
			}
			return ps
		}
		mas.addCloseList(ap.p)
		mas.addNextAPOpenList(ap)

	}
	return nil
}

//得到size最小的
func (mas *MapAStar) getMinSizeAPoint() *APoint {
	l := len(mas.openList)
	if l == 0 {
		return nil
	}
	ap := mas.openList[0]
	if l == 1 {
		return ap
	}
	for i := 1; i < l; i++ {
		if mas.openList[i].size < ap.size {
			ap = mas.openList[i]
		}
	}
	return ap
}

func (mas *MapAStar) addNextAPOpenList(ap *APoint) {
	var ap1 *APoint
	var li *line
	for l, a := range ap.a.areaMap {
		a1 := a
		if !mas.closeMap[l.sp] {
			li = &line{l.sp, ap.p}
			ap1 = &APoint{
				a:        a1,
				p:        l.sp,
				size:     ap.size + li.Distance(),
				parentAP: ap,
				length:   ap.length + 1,
			}
			mas.openList = append(mas.openList, ap1)
		}
		if !mas.closeMap[l.ep] {
			li = &line{l.ep, ap.p}
			ap1 = &APoint{
				a:        a1,
				p:        l.ep,
				size:     ap.size + li.Distance(),
				parentAP: ap,
				length:   ap.length + 1,
			}
			mas.openList = append(mas.openList, ap1)
		}
	}

}

func (mas *MapAStar) addCloseList(p point) {
	mas.closeMap[p] = true
}
