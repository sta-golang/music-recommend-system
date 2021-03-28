package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/controller/dto"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
	"math"
	"math/rand"
	"strings"
)

type creatorService struct {
}

const (
	creatorLimit = 20

	maxSimilarNum      = 6
	maxSplitSimilarNum = 18

	creatorNotSimilarCacheFmt = "creators_%d_n"
	creatorsCacheFmt      = "creators_%d_s"
	creatorDetailCacheFmt = "creator_%d_detail"
)

var PubCreatorService = &creatorService{}

func (cs *creatorService) GetCreator(page int) ([]model.Creator, *er.Error) {
	if page <= 0 {
		page = 1
	}
	key := fmt.Sprintf(creatorsCacheFmt, page)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Creator), nil
	}
	priority := cache.Priority(math.Max(float64(cache.One), float64(cache.Ten-cache.Priority(page/2)*cache.One)))
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		creators, err := model.NewCreatorMysql().SelectCreatorsOrderBySong((page-1)*creatorLimit, creatorLimit)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if len(creators) <= 0 {
			cache.PubCacheService.Set(key, creators, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, creators, cache.Hour * int(priority) * 3, priority)
		return creators, nil
	})
	if err != nil {
		return nil, er.NewError(common.DBFindErr, err)
	}
	if ret == nil {
		return nil, nil
	}
	return ret.([]model.Creator) , nil
}

func (cs *creatorService) GetCreatorWithType(tp int, page int) ([]model.Creator, *er.Error) {
	if page <= 0 {
		page = 1
	}
	creators, err := model.NewCreatorMysql().SelectCreatorsForType(tp, (page-1)*creatorLimit, creatorLimit)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindErr, err)
	}

	return creators, nil
}

func (cs *creatorService) getCreatorWithCache(id int) (*model.Creator, error) {
	key := fmt.Sprintf(creatorNotSimilarCacheFmt, id)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.(*model.Creator), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		creator, err := model.NewCreatorMysql().SelectCreator(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if creator == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, creator, cache.Hour*24, cache.Ten)
		return creator, nil
	})
	if ret != nil {
		return ret.(*model.Creator), err
	}
	return nil, err
}

// GetCreatorDetail 获取作者详细信息
func (cs *creatorService) GetCreatorDetail(id int) (*dto.CreatorAndSimilar, *er.Error) {
	key := fmt.Sprintf(creatorDetailCacheFmt, id)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
		}
		cacheDetail := val.(*dto.CreatorAndSimilar)
		if len(cacheDetail.SimilarCreator) <= maxSimilarNum {
			return cacheDetail, nil
		}
		// copy一份数据
		tempDetail := *cacheDetail
		// 如果一个作者的相似作者大于最大返回值
		// 则使用洗牌算法将这数组打散。保证多次返回又不同的结果
		// 使得看起来像是在变化
		rand.Shuffle(len(tempDetail.SimilarCreator), func(i, j int) {
			tempDetail.SimilarCreator[i], tempDetail.SimilarCreator[j] = tempDetail.SimilarCreator[j], tempDetail.SimilarCreator[i]
		})
		tempDetail.SimilarCreator = tempDetail.SimilarCreator[:maxSimilarNum]
		return &tempDetail, nil
	}
	data, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		creator, err := model.NewCreatorMysql().SelectCreator(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if creator == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(fmt.Sprintf(creatorNotSimilarCacheFmt, id), creator, cache.Hour * 24, cache.Ten)
		if creator.SimilarCreator != "" {
			split := strings.Split(creator.SimilarCreator, model.CreatorDelimiter)
			similar, err := model.NewCreatorMysql().SelectCreatorForIDs(split)
			if err != nil {
				log.Error(err)
				return dto.NewCreatorAndSimilar(creator, nil), nil
			}
			return dto.NewCreatorAndSimilar(creator, similar), nil
		}
		return dto.NewCreatorAndSimilar(creator, nil), nil
	})
	if data == nil && err == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	if err != nil {
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	if data == nil {
		return nil, nil
	}
	ret := data.(*dto.CreatorAndSimilar)
	cache.PubCacheService.Set(key, ret, cache.Hour * 24, cache.Ten)

	if len(ret.SimilarCreator) > maxSimilarNum {
		tempRet := *ret
		similar := make([]model.Creator, len(ret.SimilarCreator))
		copy(similar, ret.SimilarCreator)
		rand.Shuffle(len(similar), func(i, j int) {
			similar[i], similar[j] = similar[j], similar[i]
		})
		similar = similar[:maxSimilarNum]
		tempRet.SimilarCreator = similar
		return &tempRet, nil
	}

	return ret, nil
}
