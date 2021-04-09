package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
	"github.com/valyala/bytebufferpool"
)

const (
	musicOffset = 50

	maxStatisticsDuration = time.Hour * 2
	statisticsSignal      = -1
	minStatisticsDuration = time.Second * 5

	musicDetailCacheFmt  = "music_%d_d"
	musicCreatorCacheFmt = "music_%d_%d_c"
	musicAllCacheFmt     = "music_%d_all"
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

func (ms *musicService) GetMusic(id int) (*model.Music, *er.Error) {
	key := fmt.Sprintf(musicDetailCacheFmt, id)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
		}
		return val.(*model.Music), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (i interface{}, e error) {
		music, err := model.NewMusicMysql().SelectByID(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if music == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, music, cache.Hour*24, cache.Eight)
		return music, nil
	})
	if ret == nil && err == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	if err != nil {
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	buff := bytebufferpool.Get()
	if _, err = buff.WriteString(strconv.Itoa(id)); err != nil {
		log.Error(err)
		bytebufferpool.Put(buff)
		return ret.(*model.Music), nil
	}
	PubStatisticsService.Statistics(ms.GetName(), buff, true)
	return ret.(*model.Music), nil
}

func (ms *musicService) GetMusicForCreator(creatorID, page int) ([]model.Music, *er.Error) {
	if page <= 0 {
		page = 1
	}
	key := fmt.Sprintf(musicCreatorCacheFmt, creatorID, page)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Music), nil
	}
	priority := cache.Priority(math.Max(float64(cache.One), float64(cache.Ten-cache.Priority(page/2)*cache.One)))
	ret, err := common.SingleRunGroup.Do(key, func() (i interface{}, e error) {
		musics, err := model.NewMusicMysql().SelectMusicForCreator(creatorID, (page-1)*creatorLimit, creatorLimit)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if len(musics) <= 0 {
			cache.PubCacheService.Set(key, musics, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, musics, cache.Hour*int(priority)*3, priority)
		return musics, nil
	})

	if err != nil {
		return nil, er.NewError(common.DBFindErr, err)
	}
	if ret == nil {
		return nil, nil
	}

	return ret.([]model.Music), nil
}

func (ms *musicService) GetAllMusicWithCache(ctx context.Context, page int) ([]model.Music, *er.Error) {
	if page < 1 {
		page = 1
	}
	key := fmt.Sprintf(musicAllCacheFmt, page)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Music), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		retData, rErr := ms.GetAllMusic(ctx, page)
		if rErr != nil && retData == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
		}
		cache.PubCacheService.Set(key, retData, cache.Hour*12, cache.Eight)
		return retData, rErr
	})

	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindErr, err)
	}
	if ret == nil {
		return nil, nil
	}

	return ret.([]model.Music), nil
}

func (ms *musicService) GetAllMusic(ctx context.Context, page int) ([]model.Music, error) {
	if page < 1 {
		page = 1
	}
	ret, err := model.NewMusicMysql().SelectMusicsWithPublishTime(ctx, (page-1)*musicOffset, musicOffset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return ret, nil
}

func (ms *musicService) Top51CreatorMusic(creatorID int) ([]model.Music, *er.Error) {
	return nil, nil
}

// 统计方法 统计音乐被点击了多少次
// 不需要 点一次更新一次数据库
func (ms *musicService) RegisterStatistics() {
	PubStatisticsService.Register(ms.GetName(), &StatisticsFunc{
		ParseFunc:   ms.parseStatistics,
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
