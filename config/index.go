package config

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	MySql struct {
		Url string `required:"true"`
	}
}

func LoadConfig(file string) (*Config, error) {
	var conf Config
	err := configor.New(&configor.Config{ENVPrefix: "-"}).Load(&conf, file)
	return &conf, err
}
