package plugin

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
)

const (
	maxCreatorNum     = 200
	maxCreatorPageNum = 4
)

func getMusicClickCreatorIDs(ctx context.Context, musicIDStr string) ([]int, error) {
	var ret []int
	musicID, err := strconv.Atoi(musicIDStr)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	music, sErr := service.PubMusicService.GetMusic(musicID)
	if sErr != nil && sErr.Err != nil {
		log.ErrorContext(ctx, sErr)
		return nil, sErr.Err
	}
	if music == nil {
		return nil, nil
	}
	if music.CreatorIDs == "" {
		return nil, nil
	}
	creators := strings.Split(music.CreatorIDs, model.CreatorDelimiter)
	for j := 0; j < len(creators); j++ {
		creatorID, err := strconv.Atoi(creators[j])
		if err != nil {
			log.ErrorContext(ctx, err)
			continue
		}
		ret = append(ret, creatorID)
	}
	return ret, nil
}

func CreatorRecall(request *model.FeedRequest, params string) ([]model.Music, error) {
	var ret []model.Music
	for i := range request.UserProfile.MusicClick {
		if i >= maxCreatorNum {
			break
		}
		if request.UserProfile.MusicClick[i] == "" {
			continue
		}
		creators, err := getMusicClickCreatorIDs(request.Ctx, request.UserProfile.MusicClick[i])
		if err != nil {
			log.ErrorContext(request.Ctx, err)
			continue
		}
		for _, creatorID := range creators {
			page := rand.Intn(100) % maxCreatorPageNum
			musics, err := service.PubMusicService.GetMusicForCreator(creatorID, page+1)
			if err != nil {
				log.ErrorContext(request.Ctx, err)
				continue
			}
			ret = append(ret, musics...)
		}
	}
	return ret, nil
}
