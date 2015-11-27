// main.go
package main

import (
	"log"

	"github.com/gbember/gt/navmesh"
)

func main() {

	nm, err := navmesh.NewNavMesh("../mesh.json")
	if err != nil {
		log.Fatal(err)
	}

	ps, isWalk := nm.FindPath(179, 41, 178, 886)
	log.Println(isWalk, ps)
}
