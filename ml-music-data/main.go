package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/ml-music-data/data_load"
	"github.com/sta-golang/ml-music-data/utils"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/db"
	"github.com/sta-golang/music-recommend/model"
)

const (
	applicationConfName = "music-data.yaml"
)

func main() {
	//fix()
	//TestAysncG()
	//fmt.Println(ProcessCreatorAndMusic())
	//ProcessSong()
	//ProcessD()
	//err := ProcessData()
	//fmt.Println(err)
	//myTest()
	//fmt.Println(ProcessTag())
	//fmt.Println(Download())
	music, err := model.NewMusicMysql().SelectByID(85580)
	if err != nil {
		panic(err)
	}
	fmt.Println(music)
	fmt.Println(WriterMusic2Vec())
}

func WriterMusic2Vec() error {
	file, err := os.Create("song2Vec.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	ad := data_load.NewAlgData(file)
	return ad.WriterMusic2Vec()
}

func Download() error {
	ctx := context.Background()
	if err := model.NewMusicMysql().ResetMusicProcess(ctx); err != nil {
		return err
	}
	limit := 50
	wg := sync.WaitGroup{}
	wg.Add(2)
	type chData struct {
		musics []model.Music
		notIDS []string
	}
	ch := make(chan *chData, 10)
	upload := data_load.UPLoad{}
	go func() {
		defer wg.Done()
		defer close(ch)

		for {
			dbMusics, err := model.NewMusicMysql().SelectMusicsByStatus(model.MusicDefaultStatus, 0, limit)
			if err != nil {
				log.Error(err)
				continue
			}
			if len(dbMusics) <= 0 {
				break
			}
			ids := make([]string, len(dbMusics))
			for j := 0; j < len(dbMusics); j++ {
				ids[j] = strconv.Itoa(dbMusics[j].ID)
			}
			err = model.NewMusicMysql().UpdateMusicStatusBranch(ctx, ids, model.MusicProcessStatus)
			if err != nil {
				log.Error(err)
				return
			}
			musics, notIDs, err := upload.DownloadPy(dbMusics)
			if err != nil {
				log.Error(err)
				continue
			}
			ch <- &chData{
				musics: musics,
				notIDS: notIDs,
			}
		}
	}()
	go func() {
		defer wg.Done()
		for data := range ch {
			for i := range data.musics {
				err := upload.UploadMusic(ctx, &data.musics[i])
				if err != nil {
					log.Error(err)
					continue
				}
			}
			err := model.NewMusicMysql().UpdateMusicStatusBranch(ctx, data.notIDS, model.MusicNoneHasMusicUrlStatus)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}()
	wg.Wait()
	return nil
}

func ProcessTag() error {
	crawler := data_load.WangYiYunCrawler{}
	ids, err := crawler.GetPlayListIDs()
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Println(ids)
	defer func() {
		if e := recover(); e != nil {

			log.Fatal(e)
			log.Fatal(string(debug.Stack()))
		}
	}()
	cnt := 0
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		results, err := crawler.CrawlerPlaylistsDetail(id)
		if err != nil {
			log.Error(err)
			return err
		}
		dbWriter := data_load.MysqlDataWriter{}
		err = dbWriter.LoadPlaylistForTag(results)
		if err != nil {
			log.Error(err)
			return err
		}
		cnt++
		log.ConsoleLogger.Infof("LoadPlaylistForTag finish %.2f%%", 100*float64(cnt)/float64(len(ids)))
	}

	return nil
}

func myTest() {
	str := "hello world5"
	index := len(str)
	fmt.Println(str[:index-1])
}

func fix() {
	writer := data_load.MysqlDataWriter{}
	fmt.Println(writer.FixMusic())
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
	creators, err := model.NewCreatorMysql().SelectCreatorsForStatus(0, 0, 99999)
	if err != nil {
		log.Error(err)
		return err
	}
	//log.ConsoleLogger.Infof("crawler creator ok len : %d", len(creators))
	//err = mysqlWriter.LoadCreatorToMysql(creators)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}
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
		log.ConsoleLogger.Infof("LoadMusicToMysql ok ! ")

		_, err = model.NewCreatorMysql().UpdateCreatorsForStatus(model.StatusLoadMusicFinish, creator.ID)
		if err != nil {
			log.Error(err)
			return err
		}
		log.ConsoleLogger.Info("Process creator finish")
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
	common.InitLog()
	err = utils.InitCosClient(config.GlobalConfig().CosConfig)
	if err != nil {
		panic(err)
	}
}
