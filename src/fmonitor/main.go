package main

import (
	"fmonitor/conf"
	"os"
	"flag"
	"fmonitor/flog"
	"fmt"
	"strconv"
	"os/signal"
	"syscall"
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
	/*redisConn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	ret, err:= redis.String(redisConn.Do("info"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	retArr := strings.Split(strings.TrimRight(ret, "\r\n"), "\r\n")
	retMap := make(map[string]string)
	for _,v := range retArr {
		if index := strings.Index(v, "#"); index == -1 && v != "" {
			kvArr := strings.Split(v, ":")
			retMap[kvArr[0]] = kvArr[1]
		}
	}
	redisConn.Close()
	fmt.Println(dataprovider.NewProvider(config).SaveInfoCommand("127.0.0.1:6379",retMap))*/
	for _,v := range config.Servers {
		server := v["server"]
		port, err := strconv.Atoi(v["port"])
		conntype := v["conntype"]
		if err != nil {
			flog.Fatalf(err.Error())
			os.Exit(1)
		}
		var passport string
		if v["passport"] != "" {
			passport = v["passport"]
		}
		go infoProcess(server, conntype, passport, port)
	}
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
	<-exitChan
	flog.Infof("fmonitor shut down")
}

func infoProcess(server ,conntype, passport string, port int)  {
	fmt.Println(server)
	fmt.Println(port)
	fmt.Println(conntype)
	fmt.Println(passport)
}
