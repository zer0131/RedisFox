package dataprovider

import "fmonitor/conf"

type DataProvider interface {
	SaveMemoryInfo() string
}

func GetProvider(config *conf.Config) DataProvider  {
	if config.Datatype == "sqlite" {
		return new(SqliteProvide)
	}
	return nil
}
