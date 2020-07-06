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

type Monitor struct {
	ServerId string
	server   string
	conntype string
	password string
	port     int
	probe    *util.Probe
	sqlDb    dataprovider.DataProvider
	//redisPool *redis.Pool
	maxidle     int
	maxactive   int
	idletimeout time.Duration
	redisConn   redis.Conn
}

func RunMonitor(ctx context.Context, server, conntype, password string, port int, probe *util.Probe) (*Monitor, error) {

	monitor := new(Monitor)
	monitor.server = server
	monitor.conntype = conntype
	monitor.password = password
	monitor.port = port
	monitor.maxidle = conf.ConfigVal.BaseVal.Maxidle
	monitor.maxactive = conf.ConfigVal.BaseVal.Maxactive
	monitor.idletimeout = time.Duration(conf.ConfigVal.BaseVal.Idletimeout)
	monitor.probe = probe
	monitor.ServerId = server + ":" + strconv.Itoa(port)

	rc, err := redis.Dial(monitor.conntype, monitor.ServerId, redis.DialPassword(monitor.password))
	if err != nil {
		return nil, err
	}
	monitor.redisConn = rc
	/*monitor.redisPool = &redis.Pool{
		MaxIdle:monitor.maxidle,
		MaxActive:monitor.maxactive,
		IdleTimeout:monitor.idletimeout,
	}
	if _,err := monitor.redisPool.Dial();util.CheckError(err) == false {
		return nil,err
	}*/

	sd := dataprovider.NewProvider()

	monitor.sqlDb = sd

	go monitor.loop(ctx)
	logfox.InfofWithContext(ctx, "%s monitor start", monitor.ServerId)

	return monitor, nil
}

func (m *Monitor) loop(ctx context.Context) {

LOOP:
	for {
		select {
		case <-m.probe.Chan():
			logfox.Infof("%s monitor stop", m.ServerId)
			break LOOP
		default:
			_ = m.saveRedisCommand(ctx)
		}
	}
	_ = m.redisConn.Close()
	m.probe.Done()
}

func (m *Monitor) saveRedisCommand(ctx context.Context) error {
	ret, err := redis.String(m.redisConn.Do("monitor"))
	if err != nil {
		return err
	}
	if ret != "" {
		retArr := strings.Split(ret, " ")
		if retArr == nil {
			return nil
		}
		if len(retArr) == 1 {
			return nil
		}
		var newArr []string
		if retArr[1] == "(db" || string([]byte(retArr[1])[0]) == "[" {
			newArr = append([]string{retArr[0]}, retArr[3:]...)
		}
		command := strings.ToUpper(strings.Replace(newArr[1], "\"", "", -1))
		keyName := ""
		if len(newArr) > 2 {
			keyName = strings.Replace(newArr[2], "\"", "", -1)
		}
		arguments := ""
		if len(newArr) > 3 {
			for _, v := range newArr[3:] {
				arguments += " " + strings.Replace(v, "\"", "", -1)
			}
		}
		if command != "INFO" && command != "MONITOR" {
			_ = m.sqlDb.SaveMonitorCommand(ctx, m.ServerId, command, keyName, arguments, newArr[0])
		}
	}
	return nil
}
