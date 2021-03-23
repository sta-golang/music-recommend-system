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
		return val.([]model.Creator), nil
	}
	var creators []model.Creator
	var sErr *er.Error
	priority := cache.Priority(math.Max(float64(cache.One), float64(cache.Ten-cache.Priority(page/2)*cache.One)))
	_, _ = common.SingleRunGroup.Do(key, func() (interface{}, error) {
		var err error
		creators, err = model.NewCreatorMysql().SelectCreatorsOrderBySong((page-1)*creatorLimit, creatorLimit)
		if err != nil {
			log.Error(err)
			sErr = er.NewError(common.DBFindErr, err)
			return nil, nil
		}
		cache.PubCacheService.Set(key, creators, cache.NoExpire, priority)
		return nil, nil
	})
	return creators, nil
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
	var ret *dto.CreatorAndSimilar
	var sErr *er.Error
	_, _ = common.SingleRunGroup.Do(key, func() (interface{}, error) {
		creator, err := model.NewCreatorMysql().SelectCreator(id)
		if err != nil {
			log.Error(err)
			sErr = er.NewError(common.DBFindIDErr, err)
			return nil, nil
		}
		if creator == nil {
			sErr = er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
			cache.PubCacheService.Set(key, nil, 3000, cache.One)
			return nil, nil
		}
		if creator.SimilarCreator != "" {
			split := strings.Split(creator.SimilarCreator, model.CreatorDelimiter)
			similar, err := model.NewCreatorMysql().SelectCreatorForIDs(split)
			if err != nil {
				log.Error(err)
				ret = dto.NewCreatorAndSimilar(creator, nil)
				return nil, nil
			}
			ret = dto.NewCreatorAndSimilar(creator, similar)
		}
		return nil, nil
	})
	if ret == nil && sErr != nil {
		return ret, sErr
	}
	cache.PubCacheService.Set(key, ret, cache.NoExpire, cache.Ten)

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
