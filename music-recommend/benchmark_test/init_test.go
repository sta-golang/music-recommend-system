package benchmark_test

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/db"
)

func init() {
	err := config.InitConfig("../application.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(config.GlobalConfig())
	err = db.InitDB()
	if err != nil {
		panic(err)
	}
	err = common.InitMemoryCache()
	if err != nil {
		panic(err)
	}
	common.InitLog(&config.GlobalConfig().LogConfig)
	log.ConsoleLogger.Info(config.GlobalConfig())
}