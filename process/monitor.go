package process

import (
	"RedisFox/util"
	"RedisFox/conf"
	"time"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"RedisFox/dataprovider"
	"strings"
	"github.com/zer0131/logfox"
)

type Monitor struct {
	ServerId string
	server string
	conntype string
	password string
	port int
	probe *util.Probe
	sqlDb dataprovider.DataProvider
	//redisPool *redis.Pool
	maxidle int
	maxactive int
	idletimeout time.Duration
	redisConn redis.Conn
}

func RunMonitor(server,conntype,password string, port int, config *conf.Config, probe *util.Probe) (*Monitor, error)  {

	monitor := new(Monitor)
	monitor.server = server
	monitor.conntype = conntype
	monitor.password = password
	monitor.port = port
	monitor.maxidle = config.Maxidle
	monitor.maxactive = config.Maxactive
	monitor.idletimeout = time.Duration(config.Idletimeout)
	monitor.probe = probe
	monitor.ServerId = server+":"+strconv.Itoa(port)

	rc, err:= redis.Dial(monitor.conntype,monitor.ServerId,redis.DialPassword(monitor.password))
	if util.CheckError(err) == false {
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

	sd, err := dataprovider.NewProvider(config)
	if util.CheckError(err) == false {
		monitor.redisConn.Close()
		return nil, err
	}

	monitor.sqlDb = sd

	go monitor.loop()
	logfox.Infof("%s monitor start", monitor.ServerId)

	return monitor,nil
}

func (this *Monitor) loop()  {

LOOP:
	for{
		select {
		case <- this.probe.Chan():
			logfox.Infof("%s monitor stop", this.ServerId)
			break LOOP
		default:
			this.saveRedisCommand()
		}
	}
	this.sqlDb.Close()
	this.redisConn.Close()
	this.probe.Done()
}

func (this *Monitor) saveRedisCommand() error {
	ret, err := redis.String(this.redisConn.Do("monitor"))
	if util.CheckError(err) == false {
		return err
	}
	if ret != "" {
		retArr := strings.Split(ret, " ")
		if len(retArr) == 1 {
			return nil
		}
		var newArr []string
		if retArr[1] == "(db" || string([]byte(retArr[1])[0]) == "[" {
			newArr = append([]string{retArr[0]}, retArr[3:]...)
		}
		command := strings.ToUpper(strings.Replace(newArr[1], "\"", "" , -1))
		keyName := ""
		if len(newArr) > 2 {
			keyName = strings.Replace(newArr[2], "\"", "", -1)
		}
		arguments := ""
		if len(newArr) > 3 {
			for _,v := range newArr[3:] {
				arguments += " " + strings.Replace(v, "\"", "", -1)
			}
		}
		if command != "INFO" && command != "MONITOR" {
			this.sqlDb.SaveMonitorCommand(this.ServerId, command, keyName, arguments, newArr[0])
		}
	}
	return nil
}
