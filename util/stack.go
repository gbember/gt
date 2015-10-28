// stack.go
//堆栈错误打印
package util

import (
	"fmt"
	"runtime"

	"github.com/gbember/gt/logger"
)

//用log模块记录错误调用函数栈 必须在defer中调用
func LogPanicStack() {
	if x := recover(); x != nil {
		logger.Error("%v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				logger.Error(" frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}

//用自定义函数记录错误调用函数栈 必须在defer中调用
func LogPanicStackByFun(logFun func(string)) {
	if x := recover(); x != nil {
		str := fmt.Sprintf("%v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				str = str + fmt.Sprintf(" frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
		logFun(str)
	}
}
