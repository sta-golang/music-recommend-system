package common

import (
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/config"
	"time"
)

// @see https://github.com/sta-golang/go-lib-utils/log
// InitLog 初始化日志组件
func InitLog() {
	log.ConsoleLogger = log.NewConsoleLog(log.DEBUG, "[STA:Music-Recommend-Console]")
	conf := &config.GlobalConfig().LogConfig
	logger := log.NewFileLogAndAsync(conf, time.Second*3)
	log.SetGlobalLogger(logger)
}
