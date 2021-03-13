package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/ml-music-data/data_load"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/db"
	"runtime/debug"
	"sync"
)

const (
	applicationConfName = "music-data.yaml"
)

func main() {

	err := ProcessData()
	fmt.Println(err)
}

func ProcessData() error {

	wy := data_load.WangYiYunCrawler{}
	ids, err := wy.GetPlayListIDs()
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		if e := recover(); e != nil {

			log.Fatal(e)
			log.Fatal(string(debug.Stack()))
		}
		//marshal, err := codec.API.JsonAPI.Marshal(wy.GetCreatorKeys())
		//if err != nil {
		//	log.Error(err)
		//
		//}
		//err = ioutil.WriteFile("creator.txt", marshal, 0666)
		//if err != nil {
		//	log.Error(err)

		//}
	}()
	taskNum := 0
	for _, id := range ids {
		data, err := wy.ConversionToDataWithPlaylistID(id)
		if err != nil {
			log.Error(err)
			continue
		}

		wg := sync.WaitGroup{}
		mdw := &data_load.MysqlDataWriter{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			fErr := mdw.WriterCreator(data.Creators)
			if fErr != nil {
				log.Error(err)
			}
		}()
		//err = mdw.WriterMusic(data.Musics)
		//if err != nil {
		//	log.Error(err)
		//	continue
		//}
		wg.Wait()
		//log.ConsoleLogger.Infof("taskNum : %d finish total : %d Task Completion : %.2f%%", taskNum, len(ids), float64(taskNum)/float64(len(ids)))
		//log.ConsoleLogger.Debug(system_info.MemoryUsage())
		taskNum++
	}

	return nil
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
