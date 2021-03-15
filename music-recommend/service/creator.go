package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/controller/dto"
	"github.com/sta-golang/music-recommend/model"
	"math/rand"
	"strings"
)

type creatorService struct {
}

const (
	creatorLimit = 30

	maxSimilarNum = 6
	maxSplitSimilarNum = 18
)

var PubCreatorService = &creatorService{}

func (cs *creatorService) GetCreator(page int) ([]model.Creator, *er.Error) {
	if page <= 0 {
		page = 1
	}
	creators, err := model.NewCreatorMysql().SelectCreators((page-1)*creatorLimit, creatorLimit)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindErr, err)
	}
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
	creator, err := model.NewCreatorMysql().SelectCreator(id)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	if creator == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	// 相似作者
	if creator.SimilarCreator != "" {
		split := strings.Split(creator.SimilarCreator, model.CreatorDelimiter)

		similar, err := model.NewCreatorMysql().SelectCreatorForIDs(split)
		if err != nil {
			log.Error(err)
			return dto.NewCreatorAndSimilar(creator, nil), nil
		}
		// 如果一个作者的相似作者大于最大返回值
		// 则使用洗牌算法将这数组打散。保证多次返回又不同的结果
		// 使得看起来像是在变化
		if len(similar) > maxSimilarNum {
			rand.Shuffle(len(similar), func(i, j int) {
				similar[i],similar[j] = similar[j],similar[i]
			})
			similar = similar[:maxSimilarNum]
		}
		return dto.NewCreatorAndSimilar(creator, similar), nil
	}
	return dto.NewCreatorAndSimilar(creator, nil), nil
}

