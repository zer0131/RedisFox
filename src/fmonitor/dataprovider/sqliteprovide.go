package dataprovider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmonitor/util"
	"time"
	"fmonitor/flog"
	"os"
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
	checkErr(err)
	if runSql {
		createTable()
	}
}

type SqliteProvide struct{}

func (this *SqliteProvide) SaveMemoryInfo(server string, used int, peak int) int64 {
	stmt, err := db.Prepare("INSERT INTO memory(used,peak,server,datetime) VALUES(?,?,?,?)")
	checkErr(err)
	datetime := time.Now().Format("2006-01-02 15:04:05")
	ret, err := stmt.Exec(used, peak, server, datetime)
	checkErr(err)
	id, err := ret.LastInsertId()
	checkErr(err)
	return id
}

func (this *SqliteProvide) SaveInfoCommand(server string, info map[string]string) int64  {
	return 666
}

func checkErr(err error) {
	if err != nil {
		flog.Fatalf(err.Error())
		os.Exit(1)
	}
}

func createTable() {
	sqlData := `
	CREATE TABLE IF NOT EXISTS info(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		info TEXT NOT NULL,
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
	CREATE INDEX monitor_datedime_index ON monitor (datetime DESC);
	CREATE INDEX server_index ON monitor(server ASC);
	`
	_, err := db.Exec(sqlData)
	checkErr(err)
}
