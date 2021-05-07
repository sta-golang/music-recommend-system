package plugin

import (
	"math/rand"
	"strconv"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
)

const (
	// 这里其实最好配置成参数 这里太麻烦了 就暂时不用了
	maxTagNum       = 800
	maxRecallTagNum = 5
	defRecallTagNum = 3
	maxScoreNum     = 50
	maxRandNum      = 200
)

func TagRecall(request *model.FeedRequest, params string) ([]model.Music, error) {
	var ret []model.Music
	for i := range request.UserProfile.TagIDs {
		if i > maxRecallTagNum {
			break
		}
		index := i
		if index >= defRecallTagNum && len(request.UserProfile.TagIDs) >= maxRecallTagNum {
			index = rand.Intn(len(request.UserProfile.TagIDs)-defRecallTagNum) + defRecallTagNum
		}
		tagID, err := strconv.Atoi(request.UserProfile.TagIDs[index])
		if err != nil {
			log.ErrorContext(request.Ctx, err)
			continue
		}
		musics, err := service.PubTagMusicService.GetTagMusicByTagIDWithCache(request.Ctx, tagID, maxTagNum)
		if err != nil {
			log.ErrorContext(request.Ctx, err)
			continue
		}
		if len(musics) <= 0 {
			continue
		}
		for j := 0; j < len(musics) && j < maxScoreNum; j++ {
			ret = append(ret, musics[i])
		}
		if len(musics) <= maxScoreNum {
			continue
		}
		if len(musics) <= maxScoreNum+maxRandNum {
			for j := 0; j < len(musics); j++ {
				ret = append(ret, musics[j])
			}
			continue
		}
		for j := 0; j < maxRandNum; j++ {
			index := rand.Intn(len(musics) - maxScoreNum)
			if index+maxScoreNum >= len(musics) {
				continue
			}
			ret = append(ret, musics[index+maxScoreNum])
		}
	}
	return ret, nil
}
