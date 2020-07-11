package dataprovider

import (
	"RedisFox/conf"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zer0131/logfox"
	"strconv"
	"time"
)

type SqliteProvide struct {
}

func (sp *SqliteProvide) SaveMemoryInfo(ctx context.Context, server string, used int, peak int) error {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := "insert into `memory`(`used`,`peak`,`server`,`datetime`) values(?,?,?,?)"
	err := conf.ConfigVal.MysqlServiceOrm.Exec(sqlStr, used, peak, server, datetime).Error
	if err != nil {
		logfox.ErrorfWithContext(ctx, "SaveMemoryInfo insert error:%+v", err)
		return err
	}
	return nil
}

func (sp *SqliteProvide) SaveInfoCommand(ctx context.Context, server string, info map[string]string) error {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	jsonInfo, err := json.Marshal(info)
	if err != nil {
		logfox.ErrorfWithContext(ctx, "SaveInfoCommand json error:%+v", err)
		return err
	}
	sqlStr := "insert into `info`(`server`,`info`,`datetime`) values(?,?,?)"
	err = conf.ConfigVal.MysqlServiceOrm.Exec(sqlStr, server, string(jsonInfo), datetime).Error
	if err != nil {
		logfox.ErrorfWithContext(ctx, "SaveInfoCommand insert error:%+v", err)
		return err
	}
	return nil
}

func (sp *SqliteProvide) SaveMonitorCommand(ctx context.Context, server, command, keyname, argument, timestamp string) error {
	var datetime string
	if timestamp != "" {
		timestampFloat, _ := strconv.ParseFloat(timestamp, 64)
		datetime = time.Unix(int64(timestampFloat), 0).Format("2006-01-02 15:04:05")
	} else {
		datetime = time.Now().Format("2006-01-02 15:04:05")
	}
	sqlStr := "insert into `monitor`(`server`,`command`,`arguments`,`keyname`,`datetime`) values(?,?,?,?,?)"
	err := conf.ConfigVal.MysqlServiceOrm.Exec(sqlStr, server, command, argument, keyname, datetime).Error
	if err != nil {
		logfox.ErrorfWithContext(ctx, "SaveMonitorCommand insert error:%+v", err)
		return err
	}

	return nil
}

func (sp *SqliteProvide) GetInfo(ctx context.Context, serverId string) (map[string]interface{}, error) {
	var info string
	row := conf.ConfigVal.MysqlServiceOrm.Table("info").Where("server = ?", serverId).Select("info").Order("datetime desc").Row()
	err := row.Scan(&info)
	if err != nil {
		logfox.ErrorfWithContext(ctx, "GetInfo error:%+v", err)
		return nil, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(info), &jsonMap)
	if err != nil {
		logfox.ErrorfWithContext(ctx, "GetInfo json error:%+v", err)
		return nil, err
	}
	return jsonMap, nil
}

func (sp *SqliteProvide) GetMemoryInfo(ctx context.Context, serverId, fromDate, toDate string) ([]map[string]interface{}, error) {
	rows, err := conf.ConfigVal.MysqlServiceOrm.Table("memory").Where("server=? AND datetime>=? AND datetime<=?", serverId, fromDate, toDate).Select("used,peak,datetime").Rows()
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		logfox.ErrorfWithContext(ctx, "GetMemoryInfo error:%+v", err)
		return nil, err
	}
	var ret []map[string]interface{}
	for rows.Next() {
		var (
			used     string
			peak     string
			datetime string
		)
		if err := rows.Scan(&used, &peak, &datetime); err != nil {
			logfox.Error(err.Error())
			continue
		}
		ret = append(ret, map[string]interface{}{"used": used, "peak": peak, "datetime": datetime})
	}
	return ret, nil
}

func (sp *SqliteProvide) GetCommandStats(ctx context.Context, serverId, fromDate, toDate, groupBy string) ([]map[string]interface{}, error) {
	sqlSel := "SELECT COUNT(*) AS total, strftime('%s', datetime) AS datetime FROM monitor WHERE datetime>=? AND datetime<=? AND server=? GROUP BY strftime('%s', datetime) ORDER BY datetime DESC"

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

	sqlSelFormat := fmt.Sprintf(sqlSel, queryTimeFmt, queryTimeFmt)

	rows, err := conf.ConfigVal.MysqlServiceOrm.Raw(sqlSelFormat, fromDate, toDate, serverId).Rows()
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		logfox.ErrorfWithContext(ctx, "GetCommandStats error:%+v", err)
		return nil, err
	}

	var ret []map[string]interface{}
	for rows.Next() {
		var (
			total    string
			datetime string
		)
		err := rows.Scan(&total, &datetime)
		if err != nil {
			continue
		}
		ret = append(ret, map[string]interface{}{"total": total, "datetime": datetime})
	}

	return ret, nil
}

func (sp *SqliteProvide) GetTopCommandsStats(ctx context.Context, serverId, fromDate, toDate string) ([]map[string]interface{}, error) {
	sqlSel := "SELECT command, COUNT(*) AS total FROM monitor WHERE datetime>=? AND datetime<=? AND server=? GROUP BY command ORDER BY total ASC"

	rows, err := conf.ConfigVal.MysqlServiceOrm.Raw(sqlSel, fromDate, toDate, serverId).Rows()
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		logfox.ErrorfWithContext(ctx, "GetTopCommandsStats error:%+v", err)
		return nil, err
	}

	var ret []map[string]interface{}
	for rows.Next() {
		var (
			total   string
			command string
		)
		err := rows.Scan(&command, &total) //一定要注意变量顺序
		if err != nil {
			continue
		}
		ret = append(ret, map[string]interface{}{"total": total, "command": command})
	}

	return ret, nil
}

func (sp *SqliteProvide) GetTopKeysStats(ctx context.Context, serverId, fromDate, toDate string) ([]map[string]interface{}, error) {
	sqlSel := "SELECT keyname, COUNT(*) AS total FROM monitor WHERE datetime >= ? AND datetime <= ? AND server = ? GROUP BY keyname ORDER BY total DESC LIMIT 10"

	rows, err := conf.ConfigVal.MysqlServiceOrm.Raw(sqlSel, fromDate, toDate, serverId).Rows()
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		logfox.ErrorfWithContext(ctx, "GetTopKeysStats error:%+v", err)
		return nil, err
	}

	var ret []map[string]interface{}
	for rows.Next() {
		var (
			total   string
			keyname string
		)
		err := rows.Scan(&keyname, &total) //一定要注意变量顺序
		if err != nil {
			continue
		}
		ret = append(ret, map[string]interface{}{"total": total, "keyname": keyname})
	}

	return ret, nil
}
