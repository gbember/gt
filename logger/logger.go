//日志功能
//日志文件格式 info_log_yyyy_mm_dd.log
//每日凌晨0分1秒创建当天日志文件
//设置log标准库写向当前日志文件
package logger

import (
	"fmt"
	blog "log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	log *logger

	DebugLogFun    = func(s string) { Debug(s) }
	InfoLogFun     = func(s string) { Info(s) }
	ErrorLogFun    = func(s string) { Error(s) }
	CriticalLogFun = func(s string) { Critical(s) }
)

const (
	//系统错误(关闭服务器)
	CRITICAL int = iota + 1
	//错误
	ERROR
	//信息
	INFO
	//调试
	DEBUG
)

type logger struct {
	fd       *os.File
	baseLog  *blog.Logger
	logLevel int
	logDir   string
	buf      []byte
}

func init() {
	err := StartLog("", DEBUG)
	if err != nil {
		panic(err)
	}
}

//启动日志
//dir 日志文件存放目录
//logLevel 日志等级
func StartLog(dir string, logLevel int) error {
	dir = strings.TrimSpace(dir)

	l := new(logger)
	l.logDir = dir
	l.logLevel = logLevel
	if dir == "" {
		l.baseLog = blog.New(os.Stderr, "", 0)
	} else {
		err := os.MkdirAll(log.logDir, 0660)
		if err != nil {
			return err
		}
		var logFilename string = getLogInfoFileName()
		logFile := filepath.Join(log.logDir, logFilename)
		fd, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		l.baseLog = blog.New(fd, "", 0)
	}
	log = l
	return nil
}

//调试日志
func Debug(format string, args ...interface{}) {
	if log != nil && log.logLevel >= DEBUG {
		str := fmt.Sprintf(format, args...)
		log.output("DEBUG", str)
	}
}

//信息日志
func Info(format string, args ...interface{}) {
	if log != nil && log.logLevel >= INFO {
		str := fmt.Sprintf(format, args...)
		log.output("INFO", str)
	}
}

//错误日志
func Error(format string, args ...interface{}) {
	if log != nil && log.logLevel >= ERROR {
		str := fmt.Sprintf(format, args...)
		log.output("ERROR", str)
	}
}

//系统日志
func Critical(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	log.output("CRITICAL", str)
	os.Exit(1)
}

func (l *logger) output(levelStr string, str string) {
	now := time.Now()
	pc, file, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s[%s](%s=[%s]:%d): %s\n", now.Format("======2006/01/02 15:04:05====="), levelStr,
			runtime.FuncForPC(pc).Name(), filepath.Base(file), lineno, str)
	} else {
		src = fmt.Sprintf("%s[DEBUG] %s\n", now.Format("======2006/01/02 15:04:05====="), str)
	}
	l.baseLog.Output(0, src)
}

//启动日志
//func run() error {
//	err := os.MkdirAll(log.logDir, 0660)
//	if err != nil {
//		return err
//	}
//	var logFilename string = getLogInfoFileName()
//	logFile := filepath.Join(log.logDir, logFilename)
//	fd, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		return err
//	}
//	logChan := make(chan *string, 100000)
//	log.fd = fd
//	log.logChan = logChan
//	glog.SetOutput(fd)
//	glog.SetFlags(glog.Llongfile | glog.LstdFlags)
//	glog.SetPrefix("[log] ")
//	go loop()
//	return nil
//}

//func loop() {
//	changeFDTimer()
//	defer func() {
//		if log.timer != nil {
//			log.timer.Stop()
//		}
//	}()
//	for {
//		select {
//		case s, ok := <-log.logChan:
//			if ok {
//				log.fd.WriteString(*s)
//			} else {
//				return
//			}
//		case <-log.timer.C:
//			changeFDTimer()
//			var logFilename string = getLogInfoFileName()
//			logFile := filepath.Join(log.logDir, logFilename)
//			fd, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
//			if err != nil {
//				Error("变换日志文件错误:%s", err.Error())
//			} else {
//				glog.SetOutput(fd)
//				log.fd.Close()
//				log.fd = fd
//			}

//		}
//	}
//}

//func changeFDTimer() {
//	if log.timer != nil {
//		log.timer.Stop()
//	}
//	tNow := time.Now()
//	tt := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 1, 0, time.UTC)
//	if tNow.Unix() >= tt.Unix() {
//		tt = tt.AddDate(0, 0, 1)
//	}
//	log.timer = time.NewTimer(tt.Sub(tNow))
//}

//根据日期得到log文件
func getLogInfoFileName() string {
	year, mouth, day := time.Now().Date()
	return fmt.Sprintf("info_log_%d_%d_%d.log", year, int(mouth), day)
}
