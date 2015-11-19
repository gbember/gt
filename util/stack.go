// stack.go
//堆栈错误打印
package util

import (
	"fmt"
	"runtime"

	"github.com/gbember/gt/logger"
)

//得到调用栈信息
func GetCallStack(depth int) string {
	msg := ""
	for i := 0; i < 10; i++ {
		funcName, file, line, ok := runtime.Caller(i)
		if ok {
			msg += fmt.Sprintf(" frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
		}
	}
	return msg
}

//用log模块记录错误调用函数栈 必须在defer中调用
func LogPanicStack() {
	if x := recover(); x != nil {
		logger.ErrorDepth(2, "%v", x)
		logger.Error(GetCallStack(10))
	}
}
