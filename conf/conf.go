package conf

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

var ConfigVal = &Config{}

type autoBase struct {
	Servers     []map[string]string
	Sleeptime   int
	Maxidle     int
	Maxactive   int
	Idletimeout int
	Datatype    string
	Datapath    string
	Logpath     string
	Logname     string
	Loglevel    string
	Logexpire   int
	Serverip    string
	Serverport  int
	Staticdir   string
	Tpldir      string
	Debugmode   int
}

type Config struct {
	BaseVal autoBase
	//ToDo: orm
}

func NewConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	autoBaseVal := autoBase{}
	if err := yaml.Unmarshal(data, &autoBaseVal); err != nil {
		return err
	}
	ConfigVal.BaseVal = autoBaseVal
	return nil
}
