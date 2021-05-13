// Package service provides ...
package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
)

const (
	userReadCacheKey = "userRead_%s_r"
)

type userReadService struct {
}

var PubUserReadService = &userReadService{}

func (us *userReadService) AddUserRead(ctx context.Context, username string, items []model.Item, isMemory bool) error {
	if len(items) <= 0 || username == "" {
		return nil
	}
	userRead, err := us.GetUserRead(ctx, username, isMemory)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}

	musicIDS := make([]string, 0, len(items))
	for i := range items {
		musicIDS = append(musicIDS, strconv.Itoa(items[i].Music.ID))
	}
	if userRead == nil {
		userRead = &model.DBUserRead{
			Username:  username,
			MusicRead: strings.Join(musicIDS, model.MusicReadDelimter),
		}
		cache.PubCacheService.Set(fmt.Sprintf(userReadCacheKey, username), userRead, cache.Hour*144, cache.Forever)
		if isMemory {
			return nil
		}
		err := model.NewUserReadMysql().Insert(ctx, userRead)
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		return nil
	}
	split := strings.Split(userRead.MusicRead, model.MusicReadDelimter)
	musicIDS = append(musicIDS, split...)
	musicIDS = musicIDS[:common.MinInt(len(musicIDS), 1000)]
	userRead.MusicRead = strings.Join(musicIDS, model.MusicReadDelimter)
	if isMemory {
		return nil
	}
	err = model.NewUserReadMysql().Update(ctx, userRead)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	return nil
}

func (us *userReadService) GetUserRead(ctx context.Context, username string, isMemory bool) (*model.DBUserRead, error) {
	key := fmt.Sprintf(userReadCacheKey, username)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.(*model.DBUserRead), nil
	}
	if isMemory {
		return nil, nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ret, err := model.NewUserReadMysql().SelectByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if ret == nil && !isMemory {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		if isMemory {
			cache.PubCacheService.Set(key, nil, cache.Hour*144, cache.Forever)
		} else {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.Four)
		}
		return ret, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*model.DBUserRead), nil
}
