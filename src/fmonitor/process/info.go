package process

import (
	"fmonitor/util"
	"fmonitor/flog"
	"strconv"
	"fmt"
	"fmonitor/conf"
	"time"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type Info struct {
	ServerId string
	server string
	conntype string
	passport string
	port int
	probe *util.Probe
	sleepTime int
	redisConn redis.Conn
}

func RunInfo(server,conntype,passport string, port int, config *conf.Config, probe *util.Probe) (*Info, error) {

	info := new(Info)
	info.server = server
	info.conntype = conntype
	info.passport = passport
	info.port = port
	info.probe = probe
	info.sleepTime = config.Duration
	info.ServerId = server+":"+strconv.Itoa(port)

	redisOp := redis.DialPassword(info.passport)
	rc,err := redis.Dial(info.conntype, info.ServerId, redisOp)
	if err != nil {
		return nil, err
	}
	info.redisConn = rc

	flog.Infof(info.ServerId+" start")
	go info.loop()

	return info, nil

}

func (this *Info) loop()  {
LOOP:
	for {
		select {
		case <- this.probe.Chan():
			flog.Infof(this.ServerId+" stop")
			break LOOP
		default:
			fmt.Println(this.server)
			fmt.Println(this.port)
			fmt.Println(this.conntype)
			fmt.Println(this.passport)
		}
		time.Sleep(time.Second * time.Duration(this.sleepTime))
	}
	this.redisConn.Close()
	this.probe.Done()
}

func (this *Info) saveRedisInfo()  {
	ret, err:= redis.String(this.redisConn.Do("info"))
	if err != nil {
		//
	}
	retArr := strings.Split(strings.TrimRight(ret, "\r\n"), "\r\n")
	retMap := make(map[string]string)
	for _,v := range retArr {
		if index := strings.Index(v, "#"); index == -1 && v != "" {
			kvArr := strings.Split(v, ":")
			retMap[kvArr[0]] = kvArr[1]
		}
	}
}
