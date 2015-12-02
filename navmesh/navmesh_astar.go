// navmesh_astar.go
package navmesh

import "container/heap"

type NavmeshAstar struct {
	ol     *openList      //开放列表
	cl     []bool         //关闭列表  cl     map[Point]bool //关闭列表
	srcP   Point          //起点
	srcCP  *convexPolygon //起点所在的区域
	destP  Point          //终点
	destCP *convexPolygon //终点所在的区域
	apList []astar_point
	aindex int
	*NavMesh
}

func NewNavMeshAStar(nm *NavMesh) *NavmeshAstar {
	nmastar := &NavmeshAstar{
		ol: &openList{},
		cl: make([]bool, len(nm.cache_cl), len(nm.cache_cl)),
	}
	nmastar.NavMesh = nm
	return nmastar
}

//开放列表结构
type openList []*astar_point

//A星节点结构
type astar_point struct {
	p         Point
	pindex    uint16
	cp        *convexPolygon
	size      int64
	g         int64
	h         int64
	length    int
	lineIndex int
	parentAP  *astar_point
}

func (ol openList) Len() int           { return len(ol) }
func (ol openList) Less(i, j int) bool { return ol[i].size <= ol[j].size }
func (ol openList) Swap(i, j int)      { ol[i], ol[j] = ol[j], ol[i] }
func (ol *openList) Push(x interface{}) {
	*ol = append(*ol, x.(*astar_point))
}
func (ol *openList) Pop() interface{} {
	old := *ol
	length := len(old)
	x := old[length-1]
	*ol = old[:length-1]
	return x
}

func (nmastar *NavmeshAstar) reset() {
	*nmastar.ol = (*nmastar.ol)[0:0]
	copy(nmastar.cl, nmastar.cache_cl)

	nmastar.aindex = 0
}

func (nmastar *NavmeshAstar) mallocAP() *astar_point {
	if nmastar.aindex == len(nmastar.apList) {
		for i := 0; i < 50; i++ {
			nmastar.apList = append(nmastar.apList, astar_point{})
		}
	}
	ap := &nmastar.apList[nmastar.aindex]
	nmastar.aindex++
	return ap
}

func (nmastar *NavmeshAstar) addNextAPOpenList(ap *astar_point) {
	var ap1 *astar_point
	var li line
	size := int64(0)
	if ap.cp == nmastar.destCP {
		li.sp, li.ep = ap.p, nmastar.destP
		//		size = int64(li.Distance())
		size = li.Distance2()
		ap1 = nmastar.mallocAP()
		ap1.cp = ap.cp
		ap1.p = nmastar.destP
		ap1.g = ap.g + size
		ap1.h = 0
		ap1.size = ap1.g + ap1.h
		ap1.parentAP = ap
		ap1.length = ap.length + 1

		heap.Push(nmastar.ol, ap1)
	}
	cp := ap.cp
	length := len(cp.lcs)
	var l2cp *line2CP
	for i := 0; i < length; i++ {
		l2cp = cp.lcs[i]
		if l2cp.spindex == ap.pindex || l2cp.epindex == ap.pindex {
			if ap.parentAP.cp.id != l2cp.cp.id {
				nmastar.cl[ap.pindex] = false
				ap1 = nmastar.mallocAP()
				ap1.cp = l2cp.cp
				ap1.g = ap.g
				ap1.h = ap.h
				ap1.p = ap.p
				ap1.pindex = ap.pindex
				ap1.size = ap.size
				ap1.parentAP = ap
				ap1.length = ap.length + 1

				heap.Push(nmastar.ol, ap1)
			}
		} else {
			if !nmastar.isClosed(l2cp.spindex) {
				li.sp, li.ep = nmastar.points[l2cp.spindex], ap.p
				//				size = int64(li.Distance())
				size = li.Distance2()
				li.ep = nmastar.destP
				ap1 = nmastar.mallocAP()
				ap1.cp = l2cp.cp
				ap1.p = li.sp
				ap1.pindex = l2cp.spindex
				ap1.g = ap.g + size
				//				ap1.h = int64(li.Distance())
				ap1.h = li.Distance2()
				ap1.size = ap1.g + ap1.h
				ap1.parentAP = ap
				ap1.length = ap.length + 1
				ap1.lineIndex = i

				heap.Push(nmastar.ol, ap1)
			}
			if !nmastar.isClosed(l2cp.epindex) {
				li.sp, li.ep = nmastar.points[l2cp.epindex], ap.p
				//				size = int64(li.Distance())
				size = li.Distance2()
				li.ep = nmastar.destP
				ap1 = nmastar.mallocAP()
				ap1.cp = l2cp.cp
				ap1.p = li.sp
				ap1.pindex = l2cp.epindex
				ap1.g = ap.g + size
				//				ap1.h = int64(li.Distance())
				ap1.h = li.Distance2()
				ap1.size = ap1.g + ap1.h
				ap1.parentAP = ap
				ap1.length = ap.length + 1
				ap1.lineIndex = i

				heap.Push(nmastar.ol, ap1)
			}
		}
	}

}

func (nmastar *NavmeshAstar) addCloseList(pindex uint16) {
	nmastar.cl[pindex] = true

}
func (nmastar *NavmeshAstar) isClosed(pindex uint16) bool {
	return nmastar.cl[pindex]
}

func (nmastar *NavmeshAstar) findPath() ([]Point, bool) {
	nmastar.reset()
	heap.Init(nmastar.ol)
	ap := nmastar.mallocAP()
	ap.cp = nmastar.srcCP
	ap.p = nmastar.srcP
	ap.pindex = uint16(len(nmastar.cl) - 1)
	ap.length = 1

	heap.Push(nmastar.ol, ap)

	var apx interface{}
	for nmastar.ol.Len() > 0 {

		apx = heap.Pop(nmastar.ol)
		ap = apx.(*astar_point)
		if nmastar.isClosed(ap.pindex) {
			continue
		}

		if ap.p == nmastar.destP {
			//找到路径
			ps := nmastar.get_path(ap)
			return ps, true
		}
		nmastar.addCloseList(ap.pindex)
		nmastar.addNextAPOpenList(ap)

	}
	return nil, false
}

func (nmastar *NavmeshAstar) get_path(ap *astar_point) []Point {
	//TODO 删除一些不必要的点
	tap := ap
	lines := make([]line, 0, 100)
	l := line{}
	for ap != nil {
		if ap.parentAP != nil && ap.parentAP.parentAP != nil {
			l2cp := ap.parentAP.parentAP.cp.lcs[ap.parentAP.lineIndex]
			lines = append(lines, line{sp: nmastar.points[l2cp.spindex], ep: nmastar.points[l2cp.epindex]})
			l.sp, l.ep = ap.p, ap.parentAP.parentAP.p
			isFlag := true
			for i := 0; i < len(lines); i++ {
				if !l.isIntersect(lines[i]) {
					lines = lines[0:0]
					isFlag = false
					break
				}
			}
			if isFlag {
				lines = lines[0:0]
				ap.parentAP = ap.parentAP.parentAP
				tap.length--
				continue
			}
		}
		ap = ap.parentAP
	}
	ap = tap

	ps := make([]Point, ap.length, ap.length)
	i := ap.length - 1
	for ; ap != nil; i-- {
		ps[i] = ap.p
		ap = ap.parentAP
	}
	return ps
}
