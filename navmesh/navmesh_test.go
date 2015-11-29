// navmesh_test.go
package navmesh

import "testing"

var (
	nm, _   = NewNavMesh("mesh.json")
	nmastar = NewNavMeshAStar(nm)
)

func TestNavMeshFindPath(t *testing.T) {
	ps, isWalk := nm.FindPath(nmastar, 179, 41, 178, 886)
	t.Log(isWalk, ps)
}

func BenchmarkNavMeshFindPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nm.FindPath(nmastar, 179, 41, 178, 886)
	}
}
