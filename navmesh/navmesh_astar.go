// navmesh_astar.go
package navmesh

import (
//	"log"
	"container/heap"
)

type NavmeshAstar struct {
	ol     *openList      //开放列表
	cl []Point
//	cl     map[Point]bool //关闭列表
	srcP   Point          //起点
	srcCP  *convexPolygon //起点所在的区域
	destP  Point          //终点
	destCP *convexPolygon //终点所在的区域
	apList []astar_point
	aindex int
//	num1 int
//	num2 int
}

func NewNavMeshAStar()*NavmeshAstar{
	nmastar := &NavmeshAstar{
		ol:     &openList{},
		cl:     make([]Point, 0,100),
	}
	return nmastar
}

//开放列表结构
type openList []*astar_point

//A星节点结构
type astar_point struct {
	p        Point
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

func (nmastar *NavmeshAstar)reset(){
	*nmastar.ol = (*nmastar.ol)[0:0]
	nmastar.cl = nmastar.cl[0:0]
	nmastar.aindex = 0
}

func (nmastar *NavmeshAstar)mallocAP()*astar_point{
	if nmastar.aindex==len(nmastar.apList){
		for i:=0;i<50;i++{
			nmastar.apList = append(nmastar.apList,astar_point{})
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
		ap1.length=ap.length + 1
//		ap1 = &astar_point{
//			cp:       ap.cp,
//			p:        nmastar.destP,
//			size:     li.Distance2(),
//			parentAP: ap,
//			length:   ap.length + 1,
//		}
		heap.Push(nmastar.ol, ap1)
	}
	cp := ap.cp
	length := len(cp.lcs)
	var l2cp *line2CP
	for i := 0; i < length; i++ {
		l2cp = cp.lcs[i]
//		nmastar.num2++
		if !nmastar.isClosed(l2cp.l.sp){
//		if !nmastar.cl[l2cp.l.sp] {
			if l2cp.l.sp == ap.p {
//				delete(nmastar.cl, ap.p)
				ap.cp = l2cp.cp
				heap.Push(nmastar.ol, ap)
			} else {
				li.sp, li.ep = l2cp.l.sp, nmastar.destP
				ap1 = nmastar.mallocAP()
		ap1.cp = l2cp.cp
		ap1.p = l2cp.l.sp
		ap1.size = li.Distance2()
		ap1.parentAP = ap
		ap1.length=ap.length + 1
//				ap1 = &astar_point{
//					cp:       l2cp.cp,
//					p:        l2cp.l.sp,
//					size:     li.Distance2(),
//					parentAP: ap,
//					length:   ap.length + 1,
//				}
				heap.Push(nmastar.ol, ap1)
			}
		}
//		nmastar.num2++
if !nmastar.isClosed(l2cp.l.ep){
//		if !nmastar.cl[l2cp.l.ep] {
			if l2cp.l.ep == ap.p {
//				delete(nmastar.cl, ap.p)
				ap.cp = l2cp.cp
				heap.Push(nmastar.ol, ap)
			} else {
				li.sp, li.ep = l2cp.l.ep, nmastar.destP
				ap1.cp = l2cp.cp
		ap1.p = l2cp.l.ep
		ap1.size = li.Distance2()
		ap1.parentAP = ap
		ap1.length=ap.length + 1
//				ap1 = &astar_point{
//					cp:       l2cp.cp,
//					p:        l2cp.l.ep,
//					size:     li.Distance2(),
//					parentAP: ap,
//					length:   ap.length + 1,
//				}
				heap.Push(nmastar.ol, ap1)
			}
		}
	}

}

func (nmastar *NavmeshAstar) addCloseList(p Point) {
//	nmastar.num1++
//	nmastar.cl[p] = true
	nmastar.cl = append(nmastar.cl,p)
	
}
func (nmastar *NavmeshAstar)isClosed(p Point)bool{
	for i:=0;i<len(nmastar.cl);i++{
		if p == nmastar.cl[i]{
			return true
		}
	}
	return false
}

func (nmastar *NavmeshAstar) findPath() ([]Point, bool) {
	nmastar.reset()
	heap.Init(nmastar.ol)
	ap := nmastar.mallocAP()
	ap.cp = nmastar.srcCP
	ap.p = nmastar.srcP
	ap.length = 1
//	ap := &astar_point{
//		cp:     nmastar.srcCP,
//		p:      nmastar.srcP,
//		length: 1,
//	}
	heap.Push(nmastar.ol, ap)

	var apx interface{}
	for nmastar.ol.Len() > 0 {
		apx = heap.Pop(nmastar.ol)
		ap = apx.(*astar_point)
		
//		nmastar.num2++
		if nmastar.isClosed(ap.p){
//		if nmastar.cl[ap.p] {
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
		nmastar.addCloseList(ap.p)
		nmastar.addNextAPOpenList(ap)

	}
	return nil, false
}
