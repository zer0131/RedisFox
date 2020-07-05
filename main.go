package main

import (
	"RedisFox/conf"
	"RedisFox/process"
	"RedisFox/server"
	"RedisFox/util"
	"context"
	"github.com/zer0131/logfox"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const (
	configPath = "./config/redis-fox.yaml"
)

func main() {

	//初始化配置
	err := conf.NewConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer func() {
		//如果数据库没有关闭则关闭
		if conf.ConfigVal != nil && conf.ConfigVal.MysqlServiceOrm != nil {
			_ = conf.ConfigVal.MysqlServiceOrm.Close()
		}
	}()

	//上下文
	ctx := logfox.NewContextWithLogID(context.Background())

	//log初始化
	err = logfox.Init(conf.ConfigVal.BaseVal.Logpath, conf.ConfigVal.BaseVal.Logname, conf.ConfigVal.BaseVal.Loglevel, conf.ConfigVal.BaseVal.Logexpire)
	if err != nil {
		log.Fatalf("logfox init error")
		return
	}
	//日志要最后再关闭
	defer logfox.Close()

	wg := new(sync.WaitGroup)
	closeCh := make(chan struct{})
	probe := util.NewProbe(wg, closeCh)
	defer func() {
		close(closeCh)
		wg.Wait()
	}()

	for _, v := range conf.ConfigVal.BaseVal.Servers {
		processNum := 2
		srv := v["server"]
		port, err := strconv.Atoi(v["port"])
		conntype := v["conntype"]
		if err != nil {
			logfox.ErrorfWithContext(ctx, "port convert error:%+v", err)
			return
		}
		var password string
		if v["password"] != "" {
			password = v["password"]
		}

		//开启redis info存储
		_, infoErr := process.RunInfo(ctx, srv, conntype, password, port, probe)
		if infoErr != nil {
			processNum--
		}

		//开启redis monitor
		_, monitorErr := process.RunMonitor(ctx, srv, conntype, password, port, probe)
		if monitorErr != nil {
			processNum--
		}

		if processNum > 0 {
			wg.Add(processNum)
		}
	}

	//ToDo:缺少wg
	serv := server.NewServer(ctx)
	defer serv.Stop(ctx)

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
	<-exitChan

	logfox.InfoWithContext(ctx, "RedisFox shutdown")
}
