package main

import (
	"RedisFox/conf"
	"os"
	"flag"
	"strconv"
	"os/signal"
	"syscall"
	"sync"
	"RedisFox/util"
	"RedisFox/process"
	"RedisFox/server"
	"github.com/zer0131/logfox"
	"log"
)

var cpath string
var config *conf.Config

//初始化
func init() {
	flag.StringVar(&cpath, "config", "./config/redis-fox.yaml", "config path with yml format")
	flag.Parse()
	if cpath == "" {
		log.Fatalf("config path: %s error", cpath)
		os.Exit(1)
	}
	c, err := conf.NewConfig(cpath)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	config = c
	logfox.Init(config.Logpath, config.Logname, config.Loglevel, config.Logexpire)
	StorePid("")
}

//存储pid
func StorePid(path string) {
	pid := os.Getpid()
	if len(path) == 0 {
		path = "./run_redisfox.pid"
	}

	fout, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer fout.Close()
	fout.WriteString(strconv.Itoa(pid))
}

func main() {
	//日志要最后再关闭
	defer logfox.Close()

	wg := new(sync.WaitGroup)
	closeCh := make(chan struct{})
	probe := util.NewProbe(wg,closeCh)
	defer func() {
		close(closeCh)
		wg.Wait()
	}()

	for _,v := range config.Servers {
		processNum := 2
		srv := v["server"]
		port, err := strconv.Atoi(v["port"])
		conntype := v["conntype"]
		if err != nil {
			logfox.Error(err.Error())
			os.Exit(1)
		}
		var password string
		if v["password"] != "" {
			password = v["password"]
		}

		//开启redis info存储
		_,infoErr := process.RunInfo(srv,conntype,password,port,config,probe)
		if infoErr != nil {
			processNum--
		}

		//开启redis monitor
		_,monitorErr := process.RunMonitor(srv,conntype,password,port,config,probe)
		if monitorErr != nil {
			processNum--
		}

		if processNum > 0 {
			wg.Add(processNum)
		}
	}

	//ToDo(zer):缺少wg
	serv := server.NewServer(config)
	defer serv.Stop()

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
	<-exitChan

	logfox.Info("RedisFox shutdown")
}

