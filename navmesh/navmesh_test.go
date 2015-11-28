// navmesh_test.go
package navmesh

import (
	"testing"
	"sort"
)

var (
	nm, _ = NewNavMesh("mesh.json")
	list []int
)
func init(){
	for i:=0;i<100;i++{
		list = append(list,i)
	}
	sort.Ints(list)
}
func BenchmarkSortFind(b *testing.B){
	b.Log(sort.SearchInts(list,30))
	for i:=0;i<b.N;i++{
		sort.SearchInts(list,30)
	}
}

func TestNavMeshFindPath(t *testing.T) {
	ps, isWalk := nm.FindPath(179, 41, 178, 886)
	t.Log(isWalk,ps)

}

func BenchmarkNavMeshFindPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nm.FindPath(179, 41, 178, 886)
	}
}
