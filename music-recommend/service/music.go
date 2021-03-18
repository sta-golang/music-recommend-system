package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
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


var PubMusicService = &musicService{
	statisticsChan: make(chan int, 5120),
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
	// todo
	//select {
	//case ms.statisticsChan <- music.ID:
	//default:
	//}
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
