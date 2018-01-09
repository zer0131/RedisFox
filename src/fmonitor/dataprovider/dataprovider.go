package dataprovider

import "fmonitor/conf"

type Dataprovider interface {
	SaveMemoryInfo()
}

func GetProvider(config *conf.Config) (*Dataprovider)  {
	return nil
}
