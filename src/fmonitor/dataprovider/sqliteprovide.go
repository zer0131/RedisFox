package dataprovider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmonitor/util"
	"time"
	"encoding/json"
)

var (
	dbPath string
	db     *sql.DB
)

type SqliteProvide struct{}

func NewSqliteProvide(dbPath string) (*SqliteProvide,error)  {
	runSql := true
	if isExists, _ := util.PathExists(dbPath); isExists {
		runSql = false
	}
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if util.CheckError(err) == false {
		return nil,err
	}
	if runSql {
		if cerr := createTable();cerr != nil {
			return nil,cerr
		}
	}
	return &SqliteProvide{},nil
}

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

func (this *SqliteProvide) SaveInfoCommand(server string, info map[string]string) int64 {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	jsonByte, err := json.Marshal(info)
	if util.CheckError(err) == false {
		return 0
	}
	stmt, err := db.Prepare("INSERT INTO info(server,info,datetime) VALUES(?,?,?)")
	if util.CheckError(err) == false {
		return 0
	}
	ret, err := stmt.Exec(server, string(jsonByte), datetime)
	if util.CheckError(err) == false {
		return 0
	}
	id, err := ret.LastInsertId()
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
	ret, err := stmt.Exec(server, command, argument, keyname, datetime)
	if util.CheckError(err) == false {
		return 0
	}
	id, err := ret.LastInsertId()
	if util.CheckError(err) == false {
		return 0
	}
	return id
}

/*func (this *SqliteProvide) GetInfo(server string) (map[string]interface{}, error) {
	var info string
	err := db.QueryRow("SELECT info FROM info WHERE server=? ORDER BY datetime DESC LIMIT 1", server).Scan(&info)
	if util.CheckError(err) == false {
		return nil, err
	}
	jsonMap := make(map[string]interface{})
	jsonErr := json.Unmarshal([]byte(info), &jsonMap)
	if util.CheckError(jsonErr) == false {
		return nil, jsonErr
	}
	return jsonMap, nil
}

func (this *SqliteProvide) GetMemoryInfo(server, fromDate, toDate string) ([]map[string]interface{}, error) {
	return nil, nil
}*/

func createTable() error {
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
		return err
	}
	return nil
}
