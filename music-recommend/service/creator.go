package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/dto"
	"github.com/sta-golang/music-recommend/model"
	"strings"
)

type creatorService struct {
}

const (
	creatorLimit = 30

	maxSimilarNum = 6
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
	if len(creators) <= 0 {
		return creators, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
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
	if len(creators) <= 0 {
		return creators, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	return creators, nil
}

func (cs *creatorService) GetCreatorDetail(id int) (*dto.CreatorAndSimilar, *er.Error) {
	creator, err := model.NewCreatorMysql().SelectCreator(id)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	if creator == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	if creator.SimilarCreator != "" {
		split := strings.SplitN(creator.SimilarCreator, model.CreatorDelimiter, maxSimilarNum+1)
		if len(split) > maxSimilarNum {
			split = split[:maxSimilarNum]
		}
		similar, err := model.NewCreatorMysql().SelectCreatorForIDs(split)
		if err != nil {
			log.Error(err)
			return dto.NewCreatorAndSimilar(creator, nil), nil
		}
		return dto.NewCreatorAndSimilar(creator, similar), nil
	}
	return dto.NewCreatorAndSimilar(creator, nil), nil
}
