package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"sync"
	"time"
)

const (
	musicOffset = 50

	maxStatisticsDuration = time.Hour * 2
	statisticsSignal = -1
	minStatisticsDuration = time.Second * 5
)

type musicService struct {
	statisticsChan chan int
}

var onceStatisticsService sync.Once

var PubMusicService = &musicService{
	statisticsChan: make(chan int, 5120),
}

// 统计方法 统计音乐被点击了多少次
// 不需要 点一次更新一次数据库
func (ms *musicService) statisticsService() {
	onceStatisticsService.Do(func() {
		idleTime := time.Second * 5
		ticker := time.NewTimer(idleTime)
		table := make(map[int]int32, 512)
		var tempTable map[int]int32
		db := model.NewMusicMysql()
		for  {
			select {
			case musicID := <- ms.statisticsChan:
				/**
				这使得 单线程操作table 没有并发的问题
				只有统计线程结束了 tempTable 赋值结束
				然后交给收集线程操作
				 */
				if musicID == statisticsSignal {
					for key, val := range tempTable {
						table[key] = table[key] - val
					}
					tempTable = nil
					continue
				}
				num := int32(1)
				if val, ok := table[musicID]; ok {
					num = val+1
				}
				table[musicID] = num
			case <- ticker.C:
				fn := func() {
					tempTable = make(map[int]int32, len(table))
					for key, val := range table {
						if val == 0 {
							continue
						}
						err := db.UpdateMusicHotSource(key, val)
						if err != nil {
							log.Error(err)
						}
						tempTable[key] = val
					}
					ms.statisticsChan <- statisticsSignal
				}
				err := workerpool.Submit(fn)
				if err != nil {
					go fn()
				}
				ticker.Reset(idleTime)
			}
		}
	})
}

func (ms *musicService) GetMusic(id int) (*model.Music, *er.Error) {
	music, err := model.NewMusicMysql().SelectByID(id)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	if music == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	select {
	case ms.statisticsChan <- music.ID:
	default:
	}
	return music, nil
}

func (ms *musicService) GetMusicForCreator(creatorID, page int) ([]model.Music, *er.Error) {
	if page <= 0 {
		page = 1
	}
	musics, err := model.NewMusicMysql().SelectMusicForCreator(creatorID, (page-1)*creatorLimit, creatorLimit)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindErr, err)
	}

	return musics, nil
}

func (ms *musicService) Top50CreatorMusic(creatorID int) ([]model.Music, *er.Error) {
	return nil, nil
}
