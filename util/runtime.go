package util

import (
	"bytes"
	"runtime"
)

//得到goroutine id 字符串(性能不好)
func GoroutineID() string {
	buf := make([]byte, 15)
	buf = buf[:runtime.Stack(buf, false)]
	return string(bytes.Split(buf, []byte(" "))[1])
}
