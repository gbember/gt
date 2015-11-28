// main.go
package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gbember/gt/navmesh"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	meshFileName := "../mesh.json"

	nm, err := navmesh.NewNavMesh(meshFileName)
	if err != nil {
		log.Fatal(err)
	}
	nmastar := navmesh.NewNavMeshAStar()
	ps, isWalk := nm.FindPath(nmastar,179, 41, 178, 886)
	log.Println(isWalk, ps)
	if isWalk {
		fn := "tt.cpuprof"
		f, err := os.Create(fn)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		max := int64(100000)
		st := time.Now()
		for i := int64(0); i < max; i++ {
			nm.FindPath(nmastar,179, 41, 178, 886)
		}
		nt := time.Since(st)
		log.Println(nt, nt.Nanoseconds()/max)

		pprof.StopCPUProfile()
	}
}
