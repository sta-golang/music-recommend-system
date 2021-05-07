// Package plugin provides ...
package plugin

import (
	"math/rand"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
)

const (
	maxCreatorSimilarNum = 50
)

func CreatorSimilarRecall(request *model.FeedRequest, params string) ([]model.Music, error) {
	var ret []model.Music
	clicks := request.UserProfile.MusicClick
	var bookSet *set.HashSet
	for i := 0; i < len(clicks) && i < maxCreatorSimilarNum; i++ {
		index := i
		if len(clicks) >= maxCreatorSimilarNum {
			if bookSet == nil {
				bookSet = set.NewHashSet(maxCreatorSimilarNum)
				j := 0
				for ; j < 3; j++ {
					index = rand.Intn(len(clicks)) % len(clicks)
					if !bookSet.Contains(index) {
						bookSet.Add(index)
						break
					}
				}
				if j == 3 {
					continue
				}
			}
		}
		if clicks[index] == "" {
			continue
		}
		creators, err := getMusicClickCreatorIDs(request.Ctx, clicks[index])
		if err != nil {
			log.ErrorContext(request.Ctx, err)
			continue
		}
		for i := range creators {
			creatorID := creators[i]
			creator, sErr := service.PubCreatorService.GetCreatorDetail(creatorID)
			if sErr != nil && sErr.Err != nil {
				log.ErrorContext(request.Ctx, sErr)
				continue
			}
			for _, simiCreator := range creator.SimilarCreator {
				musics, sErr := service.PubMusicService.GetMusicForCreator(simiCreator.ID, 1)
				if sErr != nil && sErr.Err != nil {
					log.ErrorContext(request.Ctx, sErr)
					continue
				}
				ret = append(ret, musics...)
			}
		}
	}
	return ret, nil
}
