package plugin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
	"github.com/sta-golang/music-recommend/service/cache"
)

const (
	rankCreatorCacheKey = "rankCreator_key_%d"
)

func CreatorRank(request *model.FeedRequest, params string) (map[int]float64, error) {
	if len(request.RecallResults) <= 0 {
		return nil, nil
	}
	creators, _ := getCreatorsWithCache(request.Ctx, 5)
	if len(creators) <= 0 {
		return nil, nil
	}
	step := -0.20
	// TODO: 使用sync.pool  <06-05-21, FOUR SEASONS> //
	ret := make(map[int]float64, len(request.RecallResults))
	for i := range request.RecallResults {
		item := request.RecallResults[i]
		if item.Music.CreatorIDs == "" {
			continue
		}
		scoreNum := 0.0
		creatorIDS := strings.Split(item.Music.CreatorIDs, model.CreatorDelimiter)
		for j := range creatorIDS {
			creatorID := creatorIDS[j]
			score := 1.0
			for k := 0; k < len(creators); k++ {
				targetFlag := false
				score = 1.0 + (float64(k) * step)
				for c := 0; c < len(creators[k]); c++ {
					if creatorID == strconv.Itoa(creators[k][c].ID) {
						targetFlag = true
						break
					}
				}
				if targetFlag {
					break
				}
			}
			scoreNum += score
		}
		ret[item.Music.ID] = scoreNum / float64(len(creatorIDS))
	}
	return ret, nil
}

func getCreatorsWithCache(ctx context.Context, pageMax int) ([][]model.Creator, error) {
	key := fmt.Sprintf(rankCreatorCacheKey, pageMax)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([][]model.Creator), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		var err error
		var ret [][]model.Creator
		for i := 1; i < pageMax; i++ {
			creators, rErr := service.PubCreatorService.GetCreator(i)
			if rErr != nil && rErr.Err != nil {
				if err == nil {
					err = rErr.Err
				}
				continue
			}
			ret = append(ret, creators)
		}
		if len(ret) <= 0 {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, err
		} else {
			cache.PubCacheService.Set(key, ret, cache.Hour*15, cache.Nine)
		}
		return ret, err
	})
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	if ret == nil {
		return nil, nil
	}
	return ret.([][]model.Creator), nil

}
