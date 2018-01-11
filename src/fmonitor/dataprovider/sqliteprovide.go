package dataprovider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmonitor/util"
	"time"
	"os"
	"encoding/json"
)

var (
	dbPath string
	db     *sql.DB
)

func init() {
	dbPath = "./data/redisfox.db"
	runSql := true
	if isExists, _ := util.PathExists(dbPath); isExists {
		runSql = false
	}
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if util.CheckError(err) == false {
		os.Exit(1)
	}
	if runSql {
		createTable()
	}
}

type SqliteProvide struct{}

func (this *SqliteProvide) SaveMemoryInfo(server string, used int, peak int) int64 {
	stmt, err := db.Prepare("INSERT INTO memory(used,peak,server,datetime) VALUES(?,?,?,?)")
	if util.CheckError(err) == false {
		return 0
	}
	datetime := time.Now().Format("2006-01-02 15:04:05")
	ret, err := stmt.Exec(used, peak, server, datetime)
	if util.CheckError(err) == false {
		return 0
	}
	id, err := ret.LastInsertId()
	if util.CheckError(err) == false {
		return 0
	}
	return id
}

func (this *SqliteProvide) SaveInfoCommand(server string, info map[string]string) int64  {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	jsonByte, err := json.Marshal(info)
	if util.CheckError(err) == false {
		return 0
	}
	jsonStr := string(jsonByte)
	stmt, err := db.Prepare("INSERT INTO info(server,info,datetime) VALUES(?,?,?)")
	if util.CheckError(err) == false {
		return 0
	}
	ret,err := stmt.Exec(server, jsonStr, datetime)
	if util.CheckError(err) == false {
		return 0
	}
	id ,err := ret.LastInsertId()
	if util.CheckError(err) == false {
		return 0
	}
	return id
}

func (this *SqliteProvide) SaveMonitorCommand(server, command, argument, keyname string) int64 {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO monitor(server,command,arguments,keyname,datetime) VALUES(?,?,?,?,?)")
	if util.CheckError(err) == false {
		return 0
	}
	ret, err := stmt.Exec(server,command,argument,keyname,datetime)
	if util.CheckError(err) == false {
		return 0
	}
	id, err := ret.LastInsertId()
	if util.CheckError(err) == false {
		return 0
	}
	return id
}

func createTable() {
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
	CREATE INDEX monitor_datedime_index ON monitor (datetime DESC);
	CREATE INDEX server_index ON monitor(server ASC);
	`
	_, err := db.Exec(sqlData)
	if util.CheckError(err) == false {
		os.Exit(1)
	}
}
