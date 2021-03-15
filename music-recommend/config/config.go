package config

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"io/ioutil"
)

type Config struct {
	ServerName string               `yaml:"server"`
	IP         string               `yaml:"ip"`
	PProf      string               `yaml:"pprof"`
	Port       string               `yaml:"port"`
	LogConfig  log.FileLogConfig    `yaml:"log"`
	DBConfigs  map[string]*DBConfig `yaml:"database"`
}

var cfg *Config

func GlobalConfig() *Config {
	return cfg
}

type DBConfig struct {
	DBName     string `yaml:"db_name"`
	UserName   string `yaml:"username`
	PassWord   string `yaml:"password"`
	Target     string `yaml:"target"`
	Args       string `yaml:"args"`
	DriverName string `yaml:"driver_name"'`
}

func (dc *DBConfig) String() string {
	psw, _ := codec.API.CryptoAPI.Encode(dc.PassWord)
	return fmt.Sprintf("DBName:%v, UserName:%v, Password:%v, Target:%v, DriverName:%v", dc.DBName, dc.UserName, psw, dc.Target, dc.DriverName)
}

func InitConfig(path string) error {
	var conf Config
	bys, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return err
	}
	err = codec.API.YamlAPI.UnMarshal(bys, &conf)
	if err != nil {
		log.Error(err)
		return err
	}
	cfg = &conf
	return nil
}