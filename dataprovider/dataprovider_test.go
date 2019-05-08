package dataprovider

import (
	"testing"
	"RedisFox/conf"
	"flag"
	"os"
	"fmt"
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
	c, err := conf.NewConfig(cpath)
	if err != nil {
		os.Exit(1)
	}
	config = c
}

func TestNewProvider(t *testing.T) {
	sqlDb,_ := NewProvider(config)
	defer sqlDb.Close()
	fmt.Println(sqlDb.GetTopCommandsStats("127.0.0.1:6379","2018-01-22 00:00:00", "2018-01-25 23:00:00"))
}
