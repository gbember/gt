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
	p        Point
	pindex   uint16
	cp       *convexPolygon
	size     int64
	length   int
	parentAP *astar_point
}

func (ol openList) Len() int           { return len(ol) }
func (ol openList) Less(i, j int) bool { return ol[i].size < ol[j].size }
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
	//	for i := 0; i < len(nmastar.cl); i++ {
	//		nmastar.cl[i] = false
	//	}
	copy(nmastar.cl, nmastar.cache_cl)
	//	length := len(nmastar.points) + 1
	//	nmastar.cl = make([]bool, length, length)
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
	if ap.cp == nmastar.destCP {
		li.sp, li.ep = ap.p, nmastar.destP
		ap1 = nmastar.mallocAP()
		ap1.cp = ap.cp
		ap1.p = nmastar.destP
		ap1.size = li.Distance2()
		ap1.parentAP = ap
		ap1.length = ap.length + 1

		heap.Push(nmastar.ol, ap1)
	}
	cp := ap.cp
	length := len(cp.lcs)
	var l2cp *line2CP
	for i := 0; i < length; i++ {
		l2cp = cp.lcs[i]
		if !nmastar.isClosed(l2cp.spindex) {
			if l2cp.spindex == ap.pindex {
				//				delete(nmastar.cl, ap.p)
				ap.cp = l2cp.cp
				heap.Push(nmastar.ol, ap)
			} else {
				li.sp, li.ep = nmastar.points[l2cp.spindex], nmastar.destP
				ap1 = nmastar.mallocAP()
				ap1.cp = l2cp.cp
				ap1.p = li.sp
				ap1.pindex = l2cp.spindex
				ap1.size = li.Distance2()
				ap1.parentAP = ap
				ap1.length = ap.length + 1

				heap.Push(nmastar.ol, ap1)
			}
		}
		if !nmastar.isClosed(l2cp.epindex) {
			if l2cp.epindex == ap.pindex {
				//				delete(nmastar.cl, ap.p)
				ap.cp = l2cp.cp
				heap.Push(nmastar.ol, ap)
			} else {
				li.sp, li.ep = nmastar.points[l2cp.epindex], nmastar.destP
				ap1.cp = l2cp.cp
				ap1.p = li.sp
				ap1.pindex = l2cp.epindex
				ap1.size = li.Distance2()
				ap1.parentAP = ap
				ap1.length = ap.length + 1

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
			ps := make([]Point, ap.length, ap.length)
			i := ap.length - 1
			for ; ap != nil; i-- {
				ps[i] = ap.p
				ap = ap.parentAP
			}
			//			log.Println(nmastar.num1,nmastar.num2)
			return ps, true
		}
		nmastar.addCloseList(ap.pindex)
		nmastar.addNextAPOpenList(ap)

	}
	return nil, false
}
