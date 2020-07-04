package dataprovider

import (
	"RedisFox/conf"
	"context"
)

type DataProvider interface {
	SaveMemoryInfo(server string, used int, peak int) int64
	SaveInfoCommand(server string, info map[string]string) int64
	SaveMonitorCommand(server, command, argument, keyname, timestamp string) int64
	Close() error
	GetInfo(serverId string) (map[string]interface{}, error)
	GetMemoryInfo(serverId, fromDate, toDate string) ([]map[string]interface{}, error)
	GetCommandStats(serverId, fromDate, toDate, groupBy string) ([]map[string]interface{}, error)
	GetTopCommandsStats(serverId, fromDate, toDate string) ([]map[string]interface{}, error)
	GetTopKeysStats(serverId, fromDate, toDate string) ([]map[string]interface{}, error)
}

func NewProvider(ctx context.Context) (DataProvider, error) {
	if conf.ConfigVal.BaseVal.Datatype == "sqlite" {
		return NewSqliteProvide(ctx)
	} else {
		return NewSqliteProvide(ctx) //默认返回
	}
}
