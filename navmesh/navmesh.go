// navmesh.go
package navmesh

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type NavMesh struct {
	gsize   int64                      //格子大小(用于定位点在某个多边形区域)
	width   int64                      //地图宽度
	heigth  int64                      //地图高度
	maxVNum int64                      //格子大小和地图宽度计算出来的横向最多格子数
	cps     []*convexPolygon           //所有多边行区域
	gcp_m   map[int64][]*convexPolygon //格子与区域关系(一个格子可能与多个区域相交) 没有找到区域的表示不能行走
}

type navMeshJson struct {
	GSize  int64     `json:"gsize"`
	Width  int64     `json:"width"`
	Heigth int64     `json:"heigth"`
	Points [][]point `json:"points"`
}

//从json数据文件新建一个*NavMesh
func NewNavMesh(meshFileName string) (*NavMesh, error) {
	data, err := ioutil.ReadFile(meshFileName)
	if err != nil {
		return nil, err
	}
	nmj := new(navMeshJson)
	err = json.Unmarshal(data, nmj)
	if err != nil {
		return nil, err
	}
	nm := new(NavMesh)
	nm.gsize = nmj.GSize
	nm.width = nmj.Width
	nm.heigth = nmj.Heigth
	nm.maxVNum = nm.width / nm.gsize
	if nm.width%nm.gsize != 0 {
		nm.maxVNum++
	}
	nm.gcp_m = make(map[int64][]*convexPolygon, 500)
	length := len(nmj.Points)
	nm.cps = make([]*convexPolygon, length, length)
	for i := 0; i < length; i++ {
		if len(nmj.Points[i]) < 3 {
			return nil, errors.New("convexPolygon point num large 3")
		}
		nm.cps[i] = newConvexPolygon(i)
		nm.cps[i].ps = nmj.Points[i]
	}
	err = nm.makeRelas()
	return nm, err
}

func (nm *NavMesh) FindPath(x1, y1, x2, y2 int64) ([]point, bool) {

	ep := point{x2, y2}
	ecp := nm.getPointCP(ep)
	if ecp == nil {
		return nil, false
	}
	sp := point{x1, y1}
	scp := nm.getPointCP(sp)

	if ecp.id == scp.id {
		return []point{sp, ep}, true
	}

	nmastar := &navmesh_astar{
		ol:     &openList{},
		cl:     make(map[point]bool, 1000),
		srcP:   sp,
		srcCP:  scp,
		destP:  ep,
		destCP: ecp,
	}
	return nmastar.findPath()
}

//是否是可走点
func (nm *NavMesh) IsWalkOfPoint(p point) bool {
	return nm.getPointCP(p) != nil
}

//得到点所在的某个区域(可能有多个  但只返回一个 返回nil表示点不在区域中)
func (nm *NavMesh) getPointCP(p point) *convexPolygon {
	gid := p.getGridNum(nm.gsize, nm.maxVNum)
	if cps, ok := nm.gcp_m[gid]; ok {
		length := len(cps)
		for i := 0; i < length; i++ {
			if cps[i].isContainPoint(p) {
				return cps[i]
			}
		}
	}
	return nil
}

//得到点所在的格子编号
func (nm *NavMesh) getGridNum(p point) int64 {
	return p.getGridNum(nm.gsize, nm.maxVNum)
}

//得到线段穿过的格子id列表
func (nm *NavMesh) getAcossGridNums(l line) []int64 {
	return l.getAcossGridNums(nm.gsize, nm.maxVNum)
}

//得到区域包含的格子id列表
func (nm *NavMesh) getGrids(cp *convexPolygon) []int64 {
	return cp.getGrids(nm.gsize, nm.maxVNum)
}

//构建关系(格子所在的区域、区域与区域相邻关系、区域不可穿过线)
func (nm *NavMesh) makeRelas() (err error) {
	cps := nm.cps
	length := len(cps)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			make_l2cp(cps[i], cps[j])
		}
	}

	var cp *convexPolygon
	gsize := nm.gsize
	maxVNum := nm.maxVNum
	for i := 0; i < length; i++ {
		cp = cps[i]
		gidList := cp.getGrids(gsize, maxVNum)

		for j := 0; j < len(gidList); j++ {
			if cps, ok := nm.gcp_m[gidList[j]]; ok {
				cps = append(cps, cp)
				nm.gcp_m[gidList[j]] = cps
			} else {
				nm.gcp_m[gidList[j]] = []*convexPolygon{cp}
			}
		}

		err = cp.makeLines()
		if err != nil {
			return
		}
	}
	return
}
