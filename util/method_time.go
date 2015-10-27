// method_time.go
package util

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	//是否记录
	isLog = false
	//记录方法(默认打印)
	logFunc = func(methodName string, nanosec int64) {
		log.Printf("log method: %s ===>> %d\n", methodName, nanosec)
	}
	nilFunc = func() {}

	//本地记录消耗时间做多的N条记录
	logRankNum    = 0
	minMethodName string
	minNanosec    int64
	mtmap         = make(map[string]int64)
	mut           = sync.Mutex{}
)

//type method_time struct {
//	methodName string
//	nanosec    int64
//}

type method_time struct {
	methodName string
	nanosec    int64
}

//设置是否记录方法运行时间
func SetMethodTimeLog(isL bool) {
	isLog = isL
}

//设置记录方法运行时间的方法(不设置默认是log打印)
func SetMethodTimeLogFun(fun func(string, int64)) {
	if fun == nil {
		isLog = false
		return
	}
	logFunc = fun
}

func SetMethodTimeLogRankNum(num int) {
	logRankNum = num
}

//记录方法运行时间(单位ns) defer调用
func DeferLogMethodTime() func() {
	if isLog {
		startTime := time.Now()
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			methodName := runtime.FuncForPC(pc).Name()
			return func() {
				nanosec := time.Since(startTime).Nanoseconds()
				logFunc(methodName, nanosec)
				update_loal_method_time_rank(methodName, nanosec)
			}
		}
	}
	return nilFunc
}

func PrintMethodTime(fun func(string)) {
	//	l := len(mtmap)
	//	ls := make([]*method_time, 0, l)
	//	for k, v := range mtmap {
	//		ls = append(ls, &method_time{methodName: k, nanosec: v})
	//	}
	//	sort.Sort(ls)
	s := fmt.Sprintf("%v\n", mtmap)
	if fun == nil {
		log.Println(s)
	} else {
		fun(s)
	}
}

//更新本地方法执行时间排行
func update_loal_method_time_rank(methodName string, nanosec int64) {
	if logRankNum > 0 {
		mut.Lock()
		if nsec, ok := mtmap[methodName]; ok {
			if nanosec > nsec {
				mtmap[methodName] = nsec
			}
		} else {
			mtmap[methodName] = nanosec
		}
		mut.Unlock()
	}
}
