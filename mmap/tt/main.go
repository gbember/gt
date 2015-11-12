// main.go
package main

import (
	"log"

	"github.com/gbember/gt/mmap"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	mmap.Test()
}
