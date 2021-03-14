package main

import (
	"flag"
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/source"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/controller"
	"github.com/sta-golang/music-recommend/db"
	"github.com/valyala/fasthttp"
	"os"
)

var (
	addr = flag.String("addr", "", "TCP address to listen to")
)

const (
	AddrEnv = "introductionAddr"
	proEnv  = "proEnv"
	defAddr = ":8080"
)

func main() {
	flag.Parse()
	defer func() {
		source.Sync()
		if er := recover(); er != nil {
			panic(er)
		}
	}()
	if *addr == "" {
		*addr = os.Getenv(AddrEnv)
	}
	if *addr == "" && config.GlobalConfig().Port != "" {
		*addr = fmt.Sprintf("%s:%s", config.GlobalConfig().IP, config.GlobalConfig().Port)
	}
	if *addr == "" {
		*addr = defAddr
	}
	log.Info("init addr : ", *addr)
	router := controller.GlobalRouter()
	log.Fatal(fasthttp.ListenAndServe(*addr, router.Handler))
}

func init() {
	err := config.InitConfig("application.yaml")
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
