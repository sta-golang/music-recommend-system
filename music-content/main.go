package main

import (
	"context"

	"github.com/sta-golang/music-content/service"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/db"
)

func main() {
	service.PubTagMusicService.ReLoadTagMusic(context.Background())
}

func init() {
	err := config.InitConfig("../conf/application.yaml")
	if err != nil {
		panic(err)
	}
	err = db.InitDB()
	if err != nil {
		panic(err)
	}
}
