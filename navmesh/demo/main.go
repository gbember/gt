// main.go
package main

import (
	"time"
	"log"
"os"
	"github.com/gbember/gt/navmesh"
	"runtime/pprof"
)

func main() {
	log.SetFlags(log.Lshortfile|log.Ldate)
	nm, err := navmesh.NewNavMesh("../mesh.json")
	if err != nil {
		log.Fatal(err)
	}

	ps, isWalk := nm.FindPath(179, 41, 178, 886)
	log.Println(isWalk, ps)
	if isWalk{
		fn := "tt.cpuprof"
		f, err := os.Create(fn)
		if err != nil{
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil{
			log.Fatal(err)
		}
		max := int64(100000)
		st := time.Now()
		for i := int64(0); i < max; i++ {
			nm.FindPath(179, 41, 178, 886)
		}
		nt := time.Since(st)
		log.Println(nt,nt.Nanoseconds()/max)
		
		pprof.StopCPUProfile()
	}
}
