package process

import (
	"fmonitor/util"
	"fmonitor/conf"
	"time"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"fmonitor/flog"
	"fmt"
	"fmonitor/dataprovider"
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
	flog.Infof(monitor.ServerId+" monitor start")

	return monitor,nil
}

func (this *Monitor) loop()  {

	this.redisConn.Send("monitor")
	this.redisConn.Flush()

LOOP:
	for{
		select {
		case m := <- this.probe.Chan():
			fmt.Println(m)
			flog.Infof(this.ServerId+" monitor stop")
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
	ret, err := redis.String(this.redisConn.Receive())
	if util.CheckError(err) == false {
		return err
	}
	fmt.Println(ret)
	return nil
}
