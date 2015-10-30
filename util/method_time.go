// method_time.go
package util

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gbember/gt/console/command"
)

var (
	//是否记录
	isLog = false
	mtmap = make(map[string]int64)
	mut   = sync.Mutex{}
)

type method_time struct {
	methodName string
	nanosec    int64
}

type MTSlice []*method_time

func (mts MTSlice) Len() int           { return len(mts) }
func (mts MTSlice) Less(i, j int) bool { return mts[i].nanosec < mts[j].nanosec }
func (mts MTSlice) Swap(i, j int)      { mts[i], mts[j] = mts[j], mts[i] }

//注册method time命令
func RegisterMTCommand() {
	isLog = true
	mt := new(commandMT)
	command.Register(mt)
}

type commandMT struct {
	buf []byte
}

func (*commandMT) Name() string {
	return "mt"
}
func (*commandMT) Help() string {
	return "method call time rank"
}

func (mt *commandMT) Run(args []string) string {
	rankNum := 10
	if len(args) != 0 {
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return mt.usage()
		}
		rankNum = i
	}
	ls := getMethodTimes(rankNum)
	mt.buf = mt.buf[:0]
	for i := 0; i < len(ls); i++ {
		s := fmt.Sprintf("\t%s\t\t\t%d\n", ls[i].methodName, ls[i].nanosec)
		mt.buf = append(mt.buf, s...)
	}
	return string(mt.buf)
}

func (*commandMT) usage() string {
	return "mt [int]"
}

//记录方法运行时间(单位ns) defer调用
func DeferMethodTime() func() {
	if isLog {
		startTime := time.Now()
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			methodName := runtime.FuncForPC(pc).Name()
			return func() {
				nanosec := time.Since(startTime).Nanoseconds()
				update_loal_method_time_rank(methodName, nanosec)
			}
		}
	}
	return nil
}

func getMethodTimes(rankNum int) []*method_time {
	l := len(mtmap)
	ls := make([]*method_time, 0, l)
	for k, v := range mtmap {
		ls = append(ls, &method_time{methodName: k, nanosec: v})
	}
	sort.Sort(MTSlice(ls))

	if rankNum > l {
		return ls
	}
	return ls[:rankNum]
}

//更新本地方法执行时间排行
func update_loal_method_time_rank(methodName string, nanosec int64) {
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
