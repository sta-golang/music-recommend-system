package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/ml-music-data/data_load"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/db"
	"github.com/sta-golang/music-recommend/model"
	"runtime/debug"
	"sync"
	"time"
)

const (
	applicationConfName = "music-data.yaml"
)

func main() {
	//TestAysncG()
	fmt.Println(ProcessCreatorAndMusic())
	//ProcessSong()
	//ProcessD()
	//err := ProcessData()
	//fmt.Println(err)
}

func TestAysncG() {
	details := make([]data_load.APIMusicDetail, 0, 1000)
	for i := 1; i <= 1000; i++ {
		details = append(details, data_load.APIMusicDetail{
			Name:        fmt.Sprintf("Test %d ", i),
			ID:          i,
			Dt:          0,
			PublishTime: 0,
			AR: []data_load.APIMusicDetailAR{
				data_load.APIMusicDetailAR{
					CreatorID:   i + 3000,
					CreatorName: "",
				},
				data_load.APIMusicDetailAR{
					CreatorID:   i + 50000,
					CreatorName: "",
				},
				data_load.APIMusicDetailAR{
					CreatorID:   i + 350000,
					CreatorName: "",
				},
			},
			AL: data_load.APIMusicDetailAL{
				TitleName: fmt.Sprintf("Test %d ", i),
				TitleUrl:  "",
			},
		})
	}
	w := data_load.MysqlDataWriter{}
	fmt.Println(w.LoadMusicToMysql(details))
}

// 执行一次 写入所有的歌手和歌曲
func ProcessCreatorAndMusic() error {
	crawler := data_load.WangYiYunCrawler{}
	mysqlWriter := data_load.MysqlDataWriter{}
	//creators, err := crawler.CrawlerAllCreatorList()
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}
	creators, err := model.NewCreatorMysql().SelectCreators(0, 99999)
	if err != nil {
		log.Error(err)
		return err
	}
	log.ConsoleLogger.Infof("crawler creator ok len : %d", len(creators))
	err = mysqlWriter.LoadCreatorToMysql(creators)
	if err != nil {
		log.Error(err)
		return err
	}
	log.ConsoleLogger.Info("LoadCreatorToMysql ok!")
	for _, creator := range creators {
		details, err := crawler.CrawlerCreatorMusic(creator.ID)
		if err != nil {
			log.Error(err)
			return err
		}
		log.ConsoleLogger.Infof("crawler creator ：%s music ok len : %d", creator.Name, len(details))
		err = mysqlWriter.LoadMusicToMysql(details)
		if err != nil {
			log.Error(err)
			return err
		}
		log.ConsoleLogger.Infof("LoadMusicToMysql ok !")
	}
	return nil
}

func ProcessSong() {
	crawler := data_load.WangYiYunCrawler{}
	music, err := crawler.CrawlerCreatorMusic(12852319)
	fmt.Println(err)
	setT := set.NewStringSet()
	for _, m := range music {
		fmt.Println(m)
		setT.Add(fmt.Sprintf("%d", m.ID))
	}
	fmt.Println(len(music))
	fmt.Println(setT.Size())
}

func ProcessD() {
	crawler := data_load.WangYiYunCrawler{}
	fmt.Println(crawler.CrawlerAllCreatorList())
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
		time.Sleep(time.Second * 5)
		err = mdw.WriterMusic(data.Musics)
		if err != nil {
			log.Error(err)
			continue
		}
		wg.Wait()
		log.ConsoleLogger.Infof("taskNum : %d finish total : %d Task Completion : %.2f%%", taskNum, len(ids), 100*float64(taskNum)/float64(len(ids)))
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
