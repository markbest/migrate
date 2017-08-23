package conf

import (
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

var (
	Conf              config
	defaultConfigFile = "./conf/conf.toml"
)

type config struct {

	// 数据库配置
	DB database `toml:"database"`

	// 数据迁移配置
	Migrate migrate `toml:"migrate"`
}

type database struct {
	Host     string
	Database string
	Port     string
	User     string
	Password string
}

type migrate struct {
	Table string
}

func InitConfig() error {
	configBytes, err := ioutil.ReadFile(defaultConfigFile)
	if err != nil {
		return errors.New("config load err:" + err.Error())
	}
	_, err = toml.Decode(string(configBytes), &Conf)
	if err != nil {
		return errors.New("config decode err:" + err.Error())
	}
	return nil
}
