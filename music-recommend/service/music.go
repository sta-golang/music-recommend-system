package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/valyala/bytebufferpool"
	"strconv"
	"time"
)

const (
	musicOffset = 50

	maxStatisticsDuration = time.Hour * 2
	statisticsSignal = -1
	minStatisticsDuration = time.Second * 5
)

type musicService struct {
	table map[int]int32
}

func (ms *musicService) GetName() string {
	return "musicService"
}


var PubMusicService = &musicService{
	table: make(map[int]int32, 1024),
}

// 统计方法 统计音乐被点击了多少次
// 不需要 点一次更新一次数据库


func (ms *musicService) GetMusic(id int) (*model.Music, *er.Error) {
	music, err := model.NewMusicMysql().SelectByID(id)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	if music == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	buff := bytebufferpool.Get()
	if _, err = buff.WriteString(strconv.Itoa(id)); err != nil {
		log.Error(err)
		bytebufferpool.Put(buff)
		return music, nil
	}
	PubStatisticsService.Statistics(ms.GetName(),buff, true)
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

func (ms *musicService) RegisterStatistics() {
	PubStatisticsService.Register(ms.GetName(), &StatisticsFunc{
		ParseFunc: ms.parseStatistics ,
		ProcessFunc: ms.processStatistics,
	})
}

func (ms *musicService) processStatistics() {
	musicDB := model.NewMusicMysql()
	for key, val := range ms.table {
		if val == 0 {
			continue
		}
		err := musicDB.UpdateMusicHotSource(key, val)
		if err != nil {
			log.Error(err)
			continue
		}
		ms.table[key] = 0
	}
}
func (ms *musicService) parseStatistics(bytes []byte) {
	id, err := strconv.Atoi(str.BytesToString(bytes))
	if err != nil {
		log.Error(err)
		return
	}
	if id == 0 {
		return
	}
	cnt := int32(1)
	if val, ok := ms.table[id]; ok {
		cnt = val + int32(1)
	}
	ms.table[id] = cnt
}

