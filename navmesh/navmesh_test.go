// navmesh_test.go
package navmesh

import "testing"

var (
	nm, _ = NewNavMesh("mesh.json")
)

func TestNavMeshFindPath(t *testing.T) {
	//	cp := new(convexPolygon)
	//	cp.ps = []point{point{1, 1}, point{1, 12}, point{12, 12}, point{12, 1}}
	//	t.Log(nm.getGrids(cp))
	p := point{179, 41}
	gid := p.getGridNum(nm.gsize, nm.maxVNum)
	t.Log(gid)
	t.Log(nm.cps[0].isContainPoint(p))
	t.Log(nm.cps[0].isIntersectGrid(gid, nm.gsize, nm.maxVNum))

}

func Benchmark_NavMesh_FindPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nm.FindPath(99, 91, 139, 883)
	}
}
