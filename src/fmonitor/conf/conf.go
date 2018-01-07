package conf

import (
	"io/ioutil"
	"github.com/go-yaml/yaml"
)

type Config struct {
	Servers []map[string]string
	Duration int
}

func New(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := yaml.Unmarshal([]byte(data), config); err != nil {
		return nil, err
	}
	return config, nil
}
