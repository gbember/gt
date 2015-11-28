// navmesh_astar.go
package navmesh

import (
	"container/heap"
)

type navmesh_astar struct {
	ol     *openList      //开放列表
	cl     map[point]bool //关闭列表
	srcP   point          //起点
	srcCP  *convexPolygon //起点所在的区域
	destP  point          //终点
	destCP *convexPolygon //终点所在的区域
}

//开放列表结构
type openList []*astar_point

//A星节点结构
type astar_point struct {
	p        point
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

func (nmastar *navmesh_astar) addNextAPOpenList(ap *astar_point) {
	var ap1 *astar_point
	var li line
	if ap.cp == nmastar.destCP {
		li.sp, li.ep = ap.p, nmastar.destP
		ap1 = &astar_point{
			cp:       ap.cp,
			p:        nmastar.destP,
			size:     ap.size + int64(li.Distance()),
			parentAP: ap,
			length:   ap.length + 1,
		}
		heap.Push(nmastar.ol, ap1)
	}
	cp := ap.cp
	length := len(cp.lcs)
	var l2cp *line2CP
	for i := 0; i < length; i++ {
		l2cp = cp.lcs[i]
		if !nmastar.cl[l2cp.l.sp] {
			if l2cp.l.sp == ap.p {
				delete(nmastar.cl,ap.p)
				ap.cp = l2cp.cp
				heap.Push(nmastar.ol, ap)
			} else {
				li.sp, li.ep = l2cp.l.sp, ap.p
				ap1 = &astar_point{
					cp:       l2cp.cp,
					p:        l2cp.l.sp,
					size:     ap.size + int64(li.Distance()),
					parentAP: ap,
					length:   ap.length + 1,
				}
				heap.Push(nmastar.ol, ap1)
			}
		}
		if !nmastar.cl[l2cp.l.ep] {
			if l2cp.l.ep == ap.p {
				delete(nmastar.cl,ap.p)
				ap.cp = l2cp.cp
				heap.Push(nmastar.ol, ap)
			} else {
				li.sp, li.ep = l2cp.l.ep, ap.p
				ap1 = &astar_point{
					cp:       l2cp.cp,
					p:        l2cp.l.ep,
					size:     ap.size + int64(li.Distance()),
					parentAP: ap,
					length:   ap.length + 1,
				}
				heap.Push(nmastar.ol, ap1)
			}
		}
	}

}

func (nmastar *navmesh_astar) addCloseList(p point) {
	nmastar.cl[p] = true
}

func (nmastar *navmesh_astar) findPath() ([]point, bool) {
	heap.Init(nmastar.ol)
	ap := &astar_point{
		cp:     nmastar.srcCP,
		p:      nmastar.srcP,
		length: 1,
	}
	heap.Push(nmastar.ol, ap)

	var apx interface{}
	for nmastar.ol.Len() > 0 {
		apx = heap.Pop(nmastar.ol)
		ap = apx.(*astar_point)
		if nmastar.cl[ap.p]{
			continue
		}
		
		if ap.p == nmastar.destP {
			//找到路径
			ps := make([]point, ap.length, ap.length)
			i := ap.length - 1
			for ; ap != nil; i-- {
				ps[i] = ap.p
				ap = ap.parentAP
			}
			return ps, true
		}
		nmastar.addCloseList(ap.p)
		nmastar.addNextAPOpenList(ap)
		
	}
	return nil, false
}
