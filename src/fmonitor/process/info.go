package process

import (
	"fmonitor/util"
	"fmonitor/flog"
	"strconv"
	"fmonitor/conf"
	"time"
	"github.com/garyburd/redigo/redis"
	"strings"
	"fmonitor/dataprovider"
)

type Info struct {
	ServerId string
	server string
	conntype string
	password string
	port int
	sleepTime time.Duration
	probe *util.Probe
	redisConn redis.Conn
	sqlDb dataprovider.DataProvider
}

func RunInfo(server,conntype,password string, port int, config *conf.Config, probe *util.Probe) (*Info, error) {

	info := new(Info)
	info.server = server
	info.conntype = conntype
	info.password = password
	info.port = port
	info.probe = probe
	info.sleepTime = time.Duration(config.Sleeptime)
	info.ServerId = server+":"+strconv.Itoa(port)

	rc,err := redis.Dial(info.conntype, info.ServerId, redis.DialPassword(info.password))
	if util.CheckError(err) == false {
		return nil, err
	}

	sd, err := dataprovider.NewProvider(config)
	if util.CheckError(err) == false {
		rc.Close()
		return nil, err
	}

	info.redisConn = rc
	info.sqlDb = sd

	go info.loop()
	flog.Infof(info.ServerId+" info start")

	return info, nil

}

func (this *Info) loop()  {
LOOP:
	for {
		select {
		case <- this.probe.Chan():
			flog.Infof(this.ServerId+" info stop")
			break LOOP
		default:
			this.saveRedisInfo()
			time.Sleep(time.Second * this.sleepTime)
		}
	}
	this.saveRedisInfo()//最后执行一次info，用于退出在redis中阻塞的monitor
	this.sqlDb.Close()
	this.redisConn.Close()
	this.probe.Done()
}

func (this *Info) saveRedisInfo() error {
	ret, err:= redis.String(this.redisConn.Do("info"))
	if util.CheckError(err) == false {
		return err
	}
	retArr := strings.Split(strings.TrimRight(ret, "\r\n"), "\r\n")
	retMap := make(map[string]string)
	for _,v := range retArr {
		if index := strings.Index(v, "#"); index == -1 && v != "" {
			kvArr := strings.Split(v, ":")
			retMap[kvArr[0]] = kvArr[1]
		}
	}
	usedMemory,_ := strconv.Atoi(retMap["used_memory"])
	usedMemoryPeak,_ := strconv.Atoi(retMap["used_memory_peak"])
	this.sqlDb.SaveMemoryInfo(this.ServerId,usedMemory,usedMemoryPeak)
	this.sqlDb.SaveInfoCommand(this.ServerId, retMap)
	return nil
}
