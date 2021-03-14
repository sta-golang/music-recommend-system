package service

import (
	"fmt"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
)

const (
	musicOffset = 50
)

type musicService struct {
}

var PubMusicService = &musicService{}

func (ms *musicService) GetMusic(id int) (*model.Music, *er.Error) {
	music, err := model.NewMusicMysql().SelectByID(id)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindIDErr, err)
	}
	if music == nil {
		return nil, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	return music, nil
}

func (ms *musicService) GetMusicForCreator(creatorID, page int) ([]model.Music, *er.Error) {
	if page <= 0 {
		page = 1
	}
	musics, err := model.NewMusicMysql().SelectMusicForCreator(creatorID, (page-1)*creatorLimit, creatorLimit)
	if err != nil {
		log.Error(err)
		return nil, er.NewError(common.DBFindErr, err)
	}
	if len(musics) <= 0 {
		return musics, er.NewError(common.NotFound, fmt.Errorf(common.NotFoundMessage))
	}
	return musics, nil
}

func (ms *musicService) Top50CreatorMusic(creatorID int) ([]model.Music, *er.Error) {
	return nil, nil
}
