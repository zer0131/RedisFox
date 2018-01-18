package process

import (
	"fmonitor/util"
	"fmonitor/flog"
	"strconv"
	"fmt"
	"fmonitor/conf"
	"time"
)

type Info struct {
	server string
	conntype string
	passport string
	port int
	probe *util.Probe
	sleepTime int
}

func RunInfo(server,conntype,passport string, port int, config *conf.Config, probe *util.Probe) {
	info := &Info{
		server:   server,
		conntype: conntype,
		passport: passport,
		port:     port,
		probe:    probe,
		sleepTime:config.Duration,
	}

	go info.loop()

}

func (this *Info) loop()  {
	serInfo := this.server+":"+strconv.Itoa(this.port)
	flog.Infof(serInfo+" start")
	defer this.probe.Done()
LOOP:
	for {
		select {
		case <- this.probe.Chan():
			flog.Infof(serInfo+" shut down")
			break LOOP
		default:
			fmt.Println(this.server)
			fmt.Println(this.port)
			fmt.Println(this.conntype)
			fmt.Println(this.passport)
		}
		time.Sleep(time.Second*time.Duration(this.sleepTime))
	}
}
