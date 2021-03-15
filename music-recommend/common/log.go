package common

import (
	"github.com/sta-golang/go-lib-utils/log"
	"time"
)

// @see https://github.com/sta-golang/go-lib-utils/log
// InitLog 初始化日志组件
func InitLog(conf *log.FileLogConfig) {
	log.ConsoleLogger = log.NewConsoleLog(log.DEBUG, "[STA:Music-Recommend-Console]")
	logger := log.NewFileLogAndAsync(conf, time.Second*3)
	log.SetGlobalLogger(logger)
}
