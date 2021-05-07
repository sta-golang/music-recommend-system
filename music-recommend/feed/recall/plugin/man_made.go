package plugin

import (
	"context"
	"math/rand"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
	"github.com/sta-golang/music-recommend/service/cache"
)

const (
	manMadeCacheKey = "manMade_key"
)

func ManMadeRecall(request *model.FeedRequest, params string) ([]model.Music, error) {
	// TODO: 配置成字典定时读取  <07-05-21, FOUR SEASONS> //
	musicIDS := service.GetManMadeList()
	musics, err := getManMadeWithCache(request.Ctx, musicIDS)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return nil, err
	}
	rand.Shuffle(len(musics), func(i, j int) {
		musics[i], musics[j] = musics[j], musics[i]
	})
	return musics, nil
}

func getManMadeWithCache(ctx context.Context, ids []string) ([]model.Music, error) {
	key := manMadeCacheKey
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Music), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ret, err := model.NewMusicMysql().SelectMusicsByIDS(ctx, ids)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			cache.PubCacheService.Set(key, ret, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, ret, cache.Hour*24, cache.Ten)
		return ret, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.([]model.Music), nil
}
