package config

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/cache/memory"
	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"io/ioutil"
)

/**
配置类 通过解析yaml文件生成
这个可以长久保存 直接复制粘贴拿来用
*/
type Config struct {
	ServerName   string               `yaml:"server"`
	IP           string               `yaml:"ip"`
	PProf        string               `yaml:"pprof"`
	Port         string               `yaml:"port"`
	LogConfig    log.FileLogConfig    `yaml:"log"`
	MemoryConfig memory.Config        `yaml:"memory"`
	EmailConfig  EmailConfig          `yaml:"email"`
	DBConfigs    map[string]*DBConfig `yaml:"database"`
	CosConfig    COSConfig            `yaml:"cos"`
}

type COSConfig struct {
	URL string `yaml:"url"`
	SecretID string `yaml:"secretID"`
	SecretKey string `yaml:"secretKey"`
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

/**
这里就是用来初始化的
读取文件的字节数组
然后调用codec.API.YamlAPI.UnMarshal解析成需要的格式
*/
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
