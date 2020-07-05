package conf

import (
	"github.com/go-yaml/yaml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"time"
)

var ConfigVal = &Config{}

type autoBase struct {
	Servers             []map[string]string
	Sleeptime           int
	Maxidle             int
	Maxactive           int
	Idletimeout         int
	Datatype            string
	Datapath            string
	Datamaxopenconn     int
	Datamaxidleconn     int
	Datamaxconnlifetime int64
	Datalogmode         int
	Logpath             string
	Logname             string
	Loglevel            string
	Logexpire           int
	Serverip            string
	Serverport          int
	Staticdir           string
	Tpldir              string
	Debugmode           int
}

type Config struct {
	BaseVal         autoBase
	MysqlServiceOrm *gorm.DB
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

	//gorm初始化
	db, err := gorm.Open(autoBaseVal.Datatype, autoBaseVal.Datapath)
	if err != nil {
		return err
	}
	db.DB().SetMaxOpenConns(autoBaseVal.Datamaxopenconn)
	db.DB().SetMaxIdleConns(autoBaseVal.Datamaxidleconn)
	liftTime := time.Duration(autoBaseVal.Datamaxconnlifetime) * time.Millisecond
	db.DB().SetConnMaxLifetime(liftTime)
	logMode := false
	if autoBaseVal.Datalogmode == 1 {
		logMode = true
	}
	db.LogMode(logMode)
	if err := db.DB().Ping(); err != nil {
		return err
	}
	ConfigVal.MysqlServiceOrm = db

	err = initDBTable()
	if err != nil {
		return err
	}

	return nil
}

func initDBTable() error {
	sqlData := `
	CREATE TABLE IF NOT EXISTS info(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		info TEXT NOT NULL,
		server TEXT NOT NULL,
		datetime TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS memory(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		used INTEGER NOT NULL,
		peak INTEGER NOT NULL,
		server TEXT NOT NULL,
		datetime TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS monitor(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		arguments TEXT NOT NULL,
		keyname TEXT NOT NULL,
		server TEXT NOT NULL,
		datetime TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS keys(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expire NUMBER NOT NULL,
		persist NUMBER NOT NULL,
		server TEXT NOT NULL,
		datetime TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS monitor_datedime_index ON monitor (datetime DESC);
	CREATE INDEX IF NOT EXISTS server_index ON monitor(server ASC);
	`
	err := ConfigVal.MysqlServiceOrm.Exec(sqlData).Error
	if err != nil {
		return err
	}
	return nil
}
