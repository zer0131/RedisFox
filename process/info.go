package process

import (
	"RedisFox/conf"
	"RedisFox/dataprovider"
	"RedisFox/util"
	"context"
	"github.com/garyburd/redigo/redis"
	"github.com/zer0131/logfox"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	ServerId  string
	server    string
	conntype  string
	password  string
	port      int
	sleepTime time.Duration
	probe     *util.Probe
	redisConn redis.Conn
	sqlDb     dataprovider.DataProvider
}

func RunInfo(ctx context.Context, server, conntype, password string, port int, probe *util.Probe) (*Info, error) {

	info := new(Info)
	info.server = server
	info.conntype = conntype
	info.password = password
	info.port = port
	info.probe = probe
	info.sleepTime = time.Duration(conf.ConfigVal.BaseVal.Sleeptime)
	info.ServerId = server + ":" + strconv.Itoa(port)

	rc, err := redis.Dial(info.conntype, info.ServerId, redis.DialPassword(info.password))
	if err != nil {
		logfox.ErrorfWithContext(ctx, "redis connect error:%+v", err)
		return nil, err
	}

	sd, err := dataprovider.NewProvider(ctx)
	if util.CheckError(err) == false {
		rc.Close()
		return nil, err
	}

	info.redisConn = rc
	info.sqlDb = sd

	go info.loop()
	logfox.Infof("%s info start", info.ServerId)

	return info, nil

}

func (i *Info) loop() {
LOOP:
	for {
		select {
		case <-i.probe.Chan():
			logfox.Infof("%s info stop", i.ServerId)
			break LOOP
		default:
			i.saveRedisInfo()
			time.Sleep(time.Second * i.sleepTime)
		}
	}
	i.saveRedisInfo() //最后执行一次info，用于退出在redis中阻塞的monitor
	i.sqlDb.Close()
	i.redisConn.Close()
	i.probe.Done()
}

func (i *Info) saveRedisInfo() error {
	ret, err := redis.String(i.redisConn.Do("info"))
	if util.CheckError(err) == false {
		return err
	}
	retArr := strings.Split(strings.TrimRight(ret, "\r\n"), "\r\n")
	retMap := make(map[string]string)
	for _, v := range retArr {
		if index := strings.Index(v, "#"); index == -1 && v != "" {
			kvArr := strings.Split(v, ":")
			retMap[kvArr[0]] = kvArr[1]
		}
	}
	usedMemory, _ := strconv.Atoi(retMap["used_memory"])
	usedMemoryPeak, _ := strconv.Atoi(retMap["used_memory_peak"])
	i.sqlDb.SaveMemoryInfo(i.ServerId, usedMemory, usedMemoryPeak)
	i.sqlDb.SaveInfoCommand(i.ServerId, retMap)
	return nil
}
