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
	c, err := conf.New(cpath)
	if err != nil {
		os.Exit(1)
	}
	config = c
}

func TestNewProvider(t *testing.T) {
	fmt.Println(NewProvider(config).SaveMemoryInfo("127.0.0.1:6379", 12123, 1231231))
}
