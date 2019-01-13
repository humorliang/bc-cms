package logging

import (
	"fmt"
	"time"
	"strings"
	"com/setting"
	"log"
	"com/file"
	"os"
	"runtime"
)

//定义一个logger变量
var logger = log.New(os.Stderr, "", log.LstdFlags)

//获取日志路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s",
		setting.AppSetting.RuntimePath,
		setting.AppSetting.LogPath)
}

//获取每天的日志名文件
func getLogFileName(prefix string) string {
	return fmt.Sprintf("%s_%s.%s",
		prefix,
		strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0],
		"log",
	)
}

//debug日志
func Debug(v ...interface{}) {
	f, err := file.MustOpenSrc(getLogFileName("debug"), getLogFilePath())
	if err != nil {
		log.Fatalf("debug log open file is error:%s", err)
	}
	logger.SetOutput(f)
	setPrefix("debug")
	logger.SetFlags(log.LstdFlags)
	logger.Println(v)
}

//error日志
func Error(v ...interface{}) {
	f, err := file.MustOpenSrc(getLogFileName("error"), getLogFilePath())
	if err != nil {
		log.Fatalf("error log open file is error:%s", err)
	}
	logger.SetOutput(f)
	setPrefix("error")
	logger.SetFlags(log.LstdFlags)
	logger.Println(v)
}

//info日志
func Info(v ...interface{}) {
	f, err := file.MustOpenSrc(getLogFileName("info"), getLogFilePath())
	if err != nil {
		log.Fatalf("info log open file is error:%s", err)
	}
	logger.SetOutput(f)
	setPrefix("fatal")
	logger.SetFlags(log.LstdFlags)
	logger.Println(v)
}

//warn日志
func Warn(v ...interface{}) {
	f, err := file.MustOpenSrc(getLogFileName("warn"), getLogFilePath())
	if err != nil {
		log.Fatalf("warn log open file is error:%s", err)
	}
	logger.SetOutput(f)
	setPrefix("fatal")
	logger.SetFlags(log.LstdFlags)
	logger.Println(v)
}

//fatal日志
func Fatal(v ...interface{}) {
	f, err := file.MustOpenSrc(getLogFileName("fatal"), getLogFilePath())
	if err != nil {
		log.Fatalf("fatal log open file is error:%s", err)
	}
	logger.SetOutput(f)
	setPrefix("fatal")
	logger.SetFlags(log.LstdFlags)
	logger.Println(v)
}

func setPrefix(level string) {
	//可以返回运行时正在执行的文件名和行号 0 当前函数 1 上一层函数 2 上两层
	var prefix string
	_, filePath, line, ok := runtime.Caller(2)
	if ok {
		prefix = fmt.Sprintf("[%s][%s:%d]", level, filePath, line)
	} else {
		prefix = fmt.Sprintf("[%s]",level)
	}
	logger.SetPrefix(prefix)
}
