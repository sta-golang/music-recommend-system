package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
)

const (
	tagMusicKey = "tagMusic_%d_%d"
)

type tagMusicService struct {
}

var PubTagMusicService = &tagMusicService{}

func (tm *tagMusicService) GetTagMusicByTagID(ctx context.Context, tagID int, num int) ([]model.Music, error) {
	tmg, err := model.NewTagMusicMysql().SelectWithTagID(ctx, tagID)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	musicIDS := strings.SplitN(tmg.Musics, model.TagMusicDelimter, num+1)
	musics, err := model.NewMusicMysql().SelectMusicsByIDS(ctx, musicIDS[:common.MinInt(num, len(musicIDS))])
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return musics, nil
}

func (tm *tagMusicService) GetTagMusicByTagIDWithCache(ctx context.Context, tagID int, num int) ([]model.Music, error) {
	key := fmt.Sprintf(tagMusicKey, tagID, num)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Music), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ret, err := tm.GetTagMusicByTagID(ctx, tagID, num)
		if err != nil {
			return nil, err
		}
		if len(ret) <= 0 {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, ret, cache.Hour*12, cache.Eight)
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
