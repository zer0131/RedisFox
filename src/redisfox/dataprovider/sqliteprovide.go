package dataprovider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"redisfox/util"
	"time"
	"encoding/json"
	"strconv"
	"redisfox/flog"
	"fmt"
)


type SqliteProvide struct{
	dbPath string
	db *sql.DB
}

func NewSqliteProvide(dbPath string) (*SqliteProvide, error) {
	runSql := true
	if isExists, _ := util.PathExists(dbPath); isExists {
		runSql = false
	}
	db, err := sql.Open("sqlite3", dbPath)
	if util.CheckError(err) == false {
		return nil,err
	}
	sqliteProvide := new(SqliteProvide)
	sqliteProvide.dbPath = dbPath
	sqliteProvide.db = db
	if runSql {
		if cerr := sqliteProvide.createTable();cerr != nil {
			return nil,cerr
		}
	}
	return sqliteProvide,nil
}

func (this *SqliteProvide) SaveMemoryInfo(server string, used int, peak int) int64 {
	stmt, err := this.db.Prepare("INSERT INTO memory(used,peak,server,datetime) VALUES(?,?,?,?)")
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
	stmt, err := this.db.Prepare("INSERT INTO info(server,info,datetime) VALUES(?,?,?)")
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

func (this *SqliteProvide) SaveMonitorCommand(server, command, keyname,argument, timestamp string) int64 {
	var datetime string
	if timestamp != "" {
		timestampFloat, _ := strconv.ParseFloat(timestamp, 64)
		datetime = time.Unix(int64(timestampFloat), 0).Format("2006-01-02 15:04:05")
	} else {
		datetime = time.Now().Format("2006-01-02 15:04:05")
	}
	stmt, err := this.db.Prepare("INSERT INTO monitor(server,command,arguments,keyname,datetime) VALUES(?,?,?,?,?)")
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

func (this *SqliteProvide) GetInfo(serverId string) (map[string]interface{}, error) {
	var info string
	err := this.db.QueryRow("SELECT info FROM info WHERE server=? ORDER BY datetime DESC LIMIT 1", serverId).Scan(&info)
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

func (this *SqliteProvide) GetMemoryInfo(serverId, fromDate, toDate string) ([]map[string]interface{}, error) {
	sql := "SELECT used,peak,datetime FROM memory WHERE server=? AND datetime>=? AND datetime<=?"
	rows, err := this.db.Query(sql, serverId, fromDate, toDate)
	if util.CheckError(err) == false {
		return nil,err
	}
	var ret []map[string]interface{}
	for rows.Next() {
		var (
			used string
			peak string
			datetime string
		)
		if err := rows.Scan(&used,&peak,&datetime);err != nil {
			flog.Fatalf(err.Error())
			continue
		}
		ret = append(ret,map[string]interface{}{"used":used,"peak":peak,"datetime":datetime})
	}
	return ret, nil
}

func (this *SqliteProvide) GetCommandStats(serverId, fromDate, toDate, groupBy string) ([]map[string]interface{}, error) {
	sql := "SELECT COUNT(*) AS total, strftime('%s', datetime) AS datetime FROM monitor WHERE datetime>=? AND datetime<=? AND server=? GROUP BY strftime('%s', datetime) ORDER BY datetime DESC"

	var queryTimeFmt string
	if groupBy == "day" {
		queryTimeFmt = "%Y-%m-%d"
	} else if groupBy == "hour" {
		queryTimeFmt = "%Y-%m-%d %H"
	} else if groupBy == "minute" {
		queryTimeFmt = "%Y-%m-%d %H:%M"
	} else {
		queryTimeFmt = "%Y-%m-%d %H:%M:%S"
	}

	sql = fmt.Sprintf(sql, queryTimeFmt, queryTimeFmt)

	rows, err := this.db.Query(sql, fromDate,toDate,serverId)
	if util.CheckError(err) == false {
		return nil, err
	}

	var ret []map[string]interface{}
	for rows.Next() {
		var (
			total string
			datetime string
		)
		err := rows.Scan(&total,&datetime)
		if util.CheckError(err) == false {
			continue
		}
		ret = append(ret, map[string]interface{}{"total":total,"datetime":datetime})
	}

	return ret, nil
}

func (this *SqliteProvide) Close() error {
	return this.db.Close()
}

func (this *SqliteProvide) createTable() error {
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
	_, err := this.db.Exec(sqlData)
	if util.CheckError(err) == false {
		return err
	}
	return nil
}
