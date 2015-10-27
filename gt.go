// gt.go
package main

import (
	"github.com/gbember/gt/console"
	"github.com/gbember/gt/module"
)

func main() {
	console.RegisterConsole("127.0.0.1:12314", 10, 1024)
	module.Init()
	wait := make(chan bool)
	<-wait
}
