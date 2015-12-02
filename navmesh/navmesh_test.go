// navmesh_test.go
package navmesh

import "testing"

var (
	nm, _   = NewNavMesh("mesh.json")
	nmastar = NewNavMeshAStar(nm)
)

func TestNavMeshFindPath(t *testing.T) {
	ps, isWalk := nm.FindPath(nmastar, 179, 41, 178, 886)
	//	ps, isWalk := nm.FindPath(nmastar, 149, 409, 178, 886)
	//	ps, isWalk := nm.FindPath(nmastar, 314, 14, 331, 283)
	t.Log(isWalk, ps)
}

func BenchmarkNavMeshFindPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//		nm.FindPath(nmastar, 314, 14, 331, 283)
		nm.FindPath(nmastar, 149, 409, 178, 886)
		//		nm.FindPath(nmastar, 179, 41, 178, 886)
	}
}
