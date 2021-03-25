package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
)

type followService struct {
}

var PubFollowService = followService{}

func (fs *followService) Follow(creatorID int, username string) *er.Error {
	if _, ok := cache.PubCacheService.Get(fmt.Sprintf(creatorDetailCacheFmt, creatorID));!ok {
		return er.NewError(common.NotFound, common.CreatorNotExistErr)
	}
	err := model.NewFollowCreatorMysql().Insert(&model.FollowCreator{
		CreatorID: creatorID,
		Username:  username,
	})
	if err != nil {
		log.Error(err)
		return er.NewError(common.DBCreateErr, err)
	}
	return nil
}