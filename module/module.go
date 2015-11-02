// 模块功能支持
package module

import (
	"sync"
)

type Module interface {
	//初始化
	OnInit()
	//清理
	OnDestroy()
	//执行(closeSign chan bool用于判断是否退出)
	Run(chan bool)
}

type module struct {
	mi        Module
	closeSign chan bool
	wg        sync.WaitGroup
}

var modules []*module

//注册Module
func Register(mi Module) {
	m := new(module)
	m.mi = mi
	m.closeSign = make(chan bool, 1)

	modules = append(modules, m)
}

func Init() {
	for i := 0; i < len(modules); i++ {
		modules[i].mi.OnInit()
	}

	for i := 0; i < len(modules); i++ {
		go run(modules[i])
	}
}

func Destroy() {
	if len(modules) > 0 {
		for i := len(modules) - 1; i >= 0; i-- {
			destroy(modules[i])
		}
	}
}

func run(m *module) {
	m.wg.Add(1)
	m.mi.Run(m.closeSign)
	m.wg.Done()
}

func destroy(m *module) {
	m.closeSign <- true
	m.wg.Wait()
	m.mi.OnDestroy()
}
