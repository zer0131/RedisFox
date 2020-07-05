package dataprovider

import (
	"context"
)

type DataProvider interface {
	SaveMemoryInfo(ctx context.Context, server string, used int, peak int) error
	SaveInfoCommand(ctx context.Context, server string, info map[string]string) error
	SaveMonitorCommand(ctx context.Context, server, command, argument, keyname, timestamp string) error
	GetInfo(ctx context.Context, serverId string) (map[string]interface{}, error)
	GetMemoryInfo(ctx context.Context, serverId, fromDate, toDate string) ([]map[string]interface{}, error)
	GetCommandStats(ctx context.Context, serverId, fromDate, toDate, groupBy string) ([]map[string]interface{}, error)
	GetTopCommandsStats(ctx context.Context, serverId, fromDate, toDate string) ([]map[string]interface{}, error)
	GetTopKeysStats(ctx context.Context, serverId, fromDate, toDate string) ([]map[string]interface{}, error)
}

func NewProvider(ctx context.Context) DataProvider {
	return &SqliteProvide{}
}
