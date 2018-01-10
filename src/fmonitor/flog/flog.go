package flog

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	PanicLevel int = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)


type LogFile struct {
	level    int
	logTime  int64
	fileName string
	fileFd   *os.File
}

var logFile LogFile

func init()  {
	logFile.fileName = "redisfox"
	logFile.level = 4

	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func SetLevel(level int) {
	logFile.level = level
}

func Debugf(format string, args ...interface{}) {
	if logFile.level >= DebugLevel {
		log.SetPrefix("debug ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...interface{}) {
	if logFile.level >= InfoLevel {
		log.SetPrefix("info ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if logFile.level >= WarnLevel {
		log.SetPrefix("warn ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("error ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Fatalf(format string, args ...interface{}) {
	if logFile.level >= FatalLevel {
		log.SetPrefix("fatal ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func (this LogFile) Write(buf []byte) (n int, err error) {
	if this.fileName == "" {
		fmt.Printf("consol: %s", buf)
		return len(buf), nil
	}

	if logFile.logTime+3600 < time.Now().Unix() {
		logFile.createLogFile()
		logFile.logTime = time.Now().Unix()
	}

	if logFile.fileFd == nil {
		return len(buf), nil
	}

	return logFile.fileFd.Write(buf)
}

func (this *LogFile) createLogFile() {
	logdir := "./log/"
	if index := strings.LastIndex(this.fileName, "/"); index != -1 {
		logdir = this.fileName[0:index] + "/"
		os.MkdirAll(this.fileName[0:index], os.ModePerm)
	}

	now := time.Now()
	filename := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d", this.fileName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	if err := os.Rename(this.fileName, filename); err == nil {
		go func() {
			tarCmd := exec.Command("tar", "-zcf", filename+".tar.gz", filename, "--remove-files")
			tarCmd.Run()

			rmCmd := exec.Command("/bin/sh", "-c", "find "+logdir+` -type f -mtime +2 -exec rm {} \;`)
			rmCmd.Run()
		}()
	}

	for index := 0; index < 10; index++ {
		if fd, err := os.OpenFile(this.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
			this.fileFd.Sync()
			this.fileFd.Close()
			this.fileFd = fd
			break
		}

		this.fileFd = nil
	}
}

