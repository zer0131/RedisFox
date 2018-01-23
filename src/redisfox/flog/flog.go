package flog

import (
	"fmt"
	"log"
	"os"
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

var LogObj *LogFile

type LogFile struct {
	level    int
	logTime  int64
	filePath     string
	fileName string
	fileFd   *os.File
}

func Init(fname, path string, currLevel int) {
	LogObj = &LogFile{
		level:currLevel,
		filePath:path,
		fileName:fname,
	}

	log.SetOutput(LogObj)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Debugf(format string, args ...interface{}) {
	if LogObj.level >= DebugLevel {
		log.SetPrefix("debug ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...interface{}) {
	if LogObj.level >= InfoLevel {
		log.SetPrefix("info ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if LogObj.level >= WarnLevel {
		log.SetPrefix("warn ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if LogObj.level >= ErrorLevel {
		log.SetPrefix("error ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Fatalf(format string, args ...interface{}) {
	if LogObj.level >= FatalLevel {
		log.SetPrefix("fatal ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Panicf(format string, args ...interface{}) {
	if LogObj.level >= PanicLevel {
		log.SetPrefix("panic")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func (this *LogFile) Write(buf []byte) (n int, err error) {
	if this.fileName == "" {
		fmt.Printf("consol: %s", buf)
		return len(buf), nil
	}

	if LogObj.logTime+3600 < time.Now().Unix() {
		LogObj.createLogFile()
		LogObj.logTime = time.Now().Unix()
	}

	if LogObj.fileFd == nil {
		return len(buf), nil
	}

	return LogObj.fileFd.Write(buf)
}

func (this *LogFile) createLogFile() {
	/*if index := strings.LastIndex(this.fileName, "/"); index != -1 {
		logdir = this.fileName[0:index] + "/"
		os.MkdirAll(this.fileName[0:index], os.ModePerm)
	}*/

	now := time.Now()
	filename := this.filePath + fmt.Sprintf("%s.%04d%02d%02d%02d", this.fileName, now.Year(), now.Month(), now.Day(), now.Hour())
	/*if err := os.Rename(this.fileName, filename); err == nil {
		go func() {
			tarCmd := exec.Command("tar", "-zcf", filename+".tar.gz", filename, "--remove-files")
			tarCmd.Run()

			rmCmd := exec.Command("/bin/sh", "-c", "find "+logdir+` -type f -mtime +2 -exec rm {} \;`)
			rmCmd.Run()
		}()
	}*/

	//尝试三次创建
	for index := 0; index < 3; index++ {
		if fd, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm); nil == err {
			this.fileFd.Sync()
			this.fileFd.Close()
			this.fileFd = fd
			break
		}

		this.fileFd = nil
	}
}
