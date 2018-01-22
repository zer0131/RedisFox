package main

import (
	"fmonitor/conf"
	"os"
	"flag"
	"fmonitor/flog"
	"strconv"
	"os/signal"
	"syscall"
	"sync"
	"fmonitor/util"
	"fmonitor/process"
)

var cpath string
var config *conf.Config

//初始化
func init() {
	flag.StringVar(&cpath, "config", "./conf/redis-fox.yaml", "config path with yml format")
	flag.Parse()
	if cpath == "" {
		flog.Fatalf("config path not found")
		os.Exit(1)
	}
	c, err := conf.NewConfig(cpath)
	if err != nil {
		flog.Fatalf(err.Error())
		os.Exit(1)
	}
	config = c
	flog.Init(config.Logname, config.Logpath, config.Loglevel)
	//StorePid("")
}

//存储pid
/*func StorePid(path string) {
	pid := os.Getpid()
	if len(path) == 0 {
		path = "./run_pusher.pid"
	}

	fout, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer fout.Close()
	fout.WriteString(strconv.Itoa(pid))
}*/

func main() {
	wg := new(sync.WaitGroup)
	closeCh := make(chan struct{})
	probe := util.NewProbe(wg,closeCh)
	for _,v := range config.Servers {
		processNum := 2
		server := v["server"]
		port, err := strconv.Atoi(v["port"])
		conntype := v["conntype"]
		if err != nil {
			flog.Fatalf(err.Error())
			os.Exit(1)
		}
		var password string
		if v["password"] != "" {
			password = v["password"]
		}

		//开启redis info存储
		_,infoErr := process.RunInfo(server,conntype,password,port,config,probe)
		if infoErr != nil {
			processNum--
		}

		//开启redis monitor
		_,monitorErr := process.RunMonitor(server,conntype,password,port,config,probe)
		if monitorErr != nil {
			processNum--
		}

		if processNum > 0 {
			wg.Add(processNum)
		}
	}
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
	<-exitChan
	close(closeCh)
	wg.Wait()
	flog.Infof("fmonitor shut down")
}

