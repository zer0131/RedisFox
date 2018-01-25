package dataprovider

import (
	"testing"
	"redisfox/conf"
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
	fmt.Println(sqlDb.GetMemoryInfo("127.0.0.1:6379","2018-01-22 00:00:00", "2018-01-22 01:00:00"))
}
