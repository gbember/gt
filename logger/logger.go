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
	"sync"
	"time"
)

var (
	mut sync.Mutex
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
	fd         *os.File
	baseLog    *blog.Logger
	logLevel   int
	logDir     string
	buf        []byte
	isTimerDay bool
	dayTimer   *time.Timer
	mut        sync.RWMutex
	exitCnt    chan bool
	wgExit     sync.WaitGroup
}

func init() {
	l, err := StartLog("", DEBUG, false)
	if err != nil {
		panic(err)
	}
	Export(l)
}

//启动日志
//dir 日志文件存放目录
//logLevel 日志等级
func StartLog(dir string, logLevel int, isTimerDay bool) (*logger, error) {
	dir = strings.TrimSpace(dir)
	l := new(logger)
	l.logDir = dir
	l.logLevel = logLevel
	l.isTimerDay = isTimerDay
	err := l.init()
	return l, err
}

func Export(l *logger) {
	mut.Lock()
	defer mut.Unlock()
	log = l
}

//调试日志
func Debug(format string, args ...interface{}) {
	if log != nil && log.logLevel >= DEBUG {
		str := fmt.Sprintf(format, args...)
		log.output(2, "DEBUG", str)
	}
}

//信息日志
func Info(format string, args ...interface{}) {
	if log != nil && log.logLevel >= INFO {
		str := fmt.Sprintf(format, args...)
		log.output(2, "INFO", str)
	}
}

//错误日志
func Error(format string, args ...interface{}) {
	if log != nil && log.logLevel >= ERROR {
		str := fmt.Sprintf(format, args...)
		log.output(2, "ERROR", str)
	}
}

//系统日志
func Critical(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	log.output(2, "CRITICAL", str)
	os.Exit(1)
}

//调试日志
func (l *logger) Debug(format string, args ...interface{}) {
	if l.logLevel >= DEBUG {
		str := fmt.Sprintf(format, args...)
		l.output(2, "DEBUG", str)
	}
}

//信息日志
func (l *logger) Info(format string, args ...interface{}) {
	if l.logLevel >= INFO {
		str := fmt.Sprintf(format, args...)
		l.output(2, "INFO", str)
	}
}

//错误日志
func (l *logger) Error(format string, args ...interface{}) {
	if l.logLevel >= ERROR {
		str := fmt.Sprintf(format, args...)
		l.output(2, "ERROR", str)
	}
}

//系统日志
func (l *logger) Critical(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	l.output(2, "CRITICAL", str)
	os.Exit(1)
}

func (l *logger) output(depth int, levelStr string, str string) {
	now := time.Now()
	pc, file, lineno, ok := runtime.Caller(depth)
	src := ""
	if ok {
		src = fmt.Sprintf("%s[%s](%s=[%s]:%d): %s\n", now.Format("======2006/01/02 15:04:05====="), levelStr,
			runtime.FuncForPC(pc).Name(), filepath.Base(file), lineno, str)
	} else {
		src = fmt.Sprintf("%s[DEBUG] %s\n", now.Format("======2006/01/02 15:04:05====="), str)
	}
	l.mut.RLock()
	l.baseLog.Output(0, src)
	l.mut.RUnlock()
}

func (l *logger) init() error {
	if l.logDir == "" {
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
		l.fd = fd
		if l.isTimerDay {
			l.exitCnt = make(chan bool)
			l.wgExit.Add(1)
			go l.loop()
		}
	}
	return nil
}

func (l *logger) timer_day() {
	now := time.Now()
	d := time.Hour*24 -
		time.Duration(now.Hour())*time.Hour -
		time.Duration(now.Minute())*time.Minute -
		time.Duration(now.Second())*time.Second + time.Second
	l.dayTimer = time.NewTimer(d)
}

func (l *logger) Stop() {
	if l.dayTimer != nil {
		l.dayTimer.Stop()
		close(l.exitCnt)
		l.wgExit.Wait()
	}
	l.fd.Close()
}

func (l *logger) loop() {
	defer l.wgExit.Done()
	l.timer_day()
	for {
		select {
		case <-l.dayTimer.C:
			var logFilename string = getLogInfoFileName()
			logFile := filepath.Join(log.logDir, logFilename)
			fd, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				Error("day timer change log file error:%v", err)
				l.timer_day()
				continue
			}
			l.mut.Lock()
			l.fd.Close()
			l.baseLog = blog.New(fd, "", 0)
			l.fd = fd
			l.mut.Unlock()
			l.timer_day()
		case <-l.exitCnt:
			break
		}
	}
}

//根据日期得到log文件
func getLogInfoFileName() string {
	year, mouth, day := time.Now().Date()
	return fmt.Sprintf("info_log_%d_%d_%d.log", year, int(mouth), day)
}
