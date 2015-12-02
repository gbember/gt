// navmesh.go
package navmesh

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type NavMesh struct {
	gsize    int64                      //格子大小(用于定位点在某个多边形区域)
	width    int64                      //地图宽度
	heigth   int64                      //地图高度
	maxVNum  int64                      //格子大小和地图宽度计算出来的横向最多格子数
	cps      []*convexPolygon           //所有多边行区域
	gcp_m    map[int64][]*convexPolygon //格子与区域关系(一个格子可能与多个区域相交) 没有找到区域的表示不能行走
	points   []Point                    //每个顶点一个index
	cache_cl []bool
}

type NavMeshJson struct {
	GSize  int64     `json:"gsize"`
	Width  int64     `json:"width"`
	Heigth int64     `json:"heigth"`
	Points [][]Point `json:"points"`
}

//从json数据文件新建一个*NavMesh
func NewNavMesh(meshFileName string) (*NavMesh, error) {
	data, err := ioutil.ReadFile(meshFileName)
	if err != nil {
		return nil, err
	}
	nmj := new(NavMeshJson)
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
	index := uint16(0)
	pindex_m := make(map[Point]uint16)
	for i := 0; i < length; i++ {
		if len(nmj.Points[i]) < 3 {
			return nil, errors.New("convexPolygon point num large 3")
		}
		length1 := len(nmj.Points[i])
		indexs := make([]uint16, length1, length1)
		for j := 0; j < length1; j++ {
			if tindex, ok := pindex_m[nmj.Points[i][j]]; ok {
				indexs[j] = tindex
			} else {
				indexs[j] = index
				pindex_m[nmj.Points[i][j]] = index
				index++
			}
		}
		nm.cps[i] = newConvexPolygon(i)
		nm.cps[i].pindexs = indexs
	}
	length = len(pindex_m)
	nm.points = make([]Point, length, length)
	for p, i := range pindex_m {
		nm.points[i] = p
	}
	err = nm.makeRelas()
	if err != nil {
		return nil, err
	}
	length = len(nm.points) + 1
	nm.cache_cl = make([]bool, length, length)

	return nm, nil
}

func (nm *NavMesh) FindPath(nmastar *NavmeshAstar, x1, y1, x2, y2 int64) ([]Point, bool) {
	ep := Point{x2, y2}
	ecp := nm.getPointCP(ep)
	if ecp == nil {
		return nil, false
	}
	sp := Point{x1, y1}
	scp := nm.getPointCP(sp)

	if ecp.id == scp.id {
		return []Point{sp, ep}, true
	}

	nmastar.srcP = sp
	nmastar.srcCP = scp
	nmastar.destP = ep
	nmastar.destCP = ecp

	return nmastar.findPath()
}

//是否是可走点
func (nm *NavMesh) IsWalkOfPoint(p Point) bool {
	return nm.getPointCP(p) != nil
}

//得到点所在的某个区域(可能有多个  但只返回一个 返回nil表示点不在区域中)
func (nm *NavMesh) getPointCP(p Point) *convexPolygon {
	gid := p.getGridNum(nm.gsize, nm.maxVNum)
	if cps, ok := nm.gcp_m[gid]; ok {
		length := len(cps)
		for i := 0; i < length; i++ {
			if cps[i].isContainPoint(nm, p) {
				return cps[i]
			}
		}
	}
	return nil
}

//得到点所在的格子编号
func (nm *NavMesh) GetGridNum(p Point) int64 {
	return p.getGridNum(nm.gsize, nm.maxVNum)
}

//得到线段穿过的格子id列表
func (nm *NavMesh) GetAcossGridNums(l line) []int64 {
	return l.getAcossGridNums(nm.gsize, nm.maxVNum)
}

//得到区域包含的格子id列表
func (nm *NavMesh) GetGrids(cp *convexPolygon) []int64 {
	return cp.getGrids(nm)
}

//构建关系(格子所在的区域、区域与区域相邻关系、区域不可穿过线)
func (nm *NavMesh) makeRelas() (err error) {
	cps := nm.cps
	length := len(cps)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			nm.make_l2cp(cps[i], cps[j])
		}
	}

	var cp *convexPolygon

	for i := 0; i < length; i++ {
		cp = cps[i]
		gidList := cp.getGrids(nm)
		for j := 0; j < len(gidList); j++ {
			if cps, ok := nm.gcp_m[gidList[j]]; ok {
				cps = append(cps, cp)
				nm.gcp_m[gidList[j]] = cps
			} else {
				nm.gcp_m[gidList[j]] = []*convexPolygon{cp}
			}
		}

		err = cp.makeLines(nm)
		if err != nil {
			return
		}
	}
	return
}
