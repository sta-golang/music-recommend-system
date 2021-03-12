package main

import (
	"fmt"
	"github.com/sta-golang/ml-music-data/data_load"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/db"
)

const (
	applicationConfName = "music-data.yaml"
)

func main() {
	wy := data_load.WangYiYunCrawler{}
	ids, err := wy.GetPlayListIDs()
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)
	data, err := wy.ConversionToDataWithPlaylistID(ids[0])
	//for _, music := range musics {
	//	err = model.NewMusicDB().InsertMusic(&music)
	//	if err != nil {
	//		log.FrameworkLogger.Error(err)
	//	}
	//}
	fmt.Println(data)
	fmt.Println(err)
	fmt.Println("-------------------------")
	for _, creator := range data.Creators {
		fmt.Println(creator)
		fmt.Println("-----------------------------------------------")
	}
}

func init() {
	err := config.InitConfig(applicationConfName)
	if err != nil {
		panic(err)
	}
	err = db.InitDB()
	if err != nil {
		panic(err)
	}
	common.InitLog(&config.GlobalConfig().LogConfig)
}
