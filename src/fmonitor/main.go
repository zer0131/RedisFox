package main

import (
	"fmonitor/conf"
	"os"
	"flag"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"strings"
	"encoding/json"
)

var cpath string
var config *conf.Config

//初始化
func init() {
	flag.StringVar(&cpath, "config", "./conf/redis-fox.yaml", "config path with yml format")
	flag.Parse()
	if cpath == "" {
		os.Exit(1)
	}
	c, err := conf.New(cpath)
	if err != nil {
		os.Exit(1)
	}
	config = c
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
	redisConn, _ := redis.Dial("tcp", "127.0.0.1:6379")
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
	retByte, _ := json.Marshal(retMap)
	fmt.Println(string(retByte))
	redisConn.Close()
	//fmt.Println(dataprovider.NewProvider(config).SaveMemoryInfo("127.0.0.1:6539", 11015696, 11321968))
}
