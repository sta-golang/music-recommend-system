package plugin

import (
	"strings"

	"github.com/sta-golang/music-recommend/model"
)

const ()

func TagRank(request *model.FeedRequest, params string) (map[int]float64, error) {
	if request.UserProfile == nil {
		return nil, nil
	}
	if len(request.UserProfile.TagIDs) <= 0 {
		return nil, nil
	}
	// TODO: 放到sync.pool中  <06-05-21, FOUR SEASONS> //
	ret := make(map[int]float64, len(request.RecallResults))
	maxScore, minScore := getTagMaxAndMin(request)
	tagMap := request.UserProfile.TagScore
	for i := range request.RecallResults {
		item := request.RecallResults[i]
		if item.Music.TagIDs == "" {
			continue
		}
		tagIDs := strings.Split(item.Music.TagIDs, model.TagDelimiter)
		scoreSum := 0.0
		cnt := 0
		for j := range tagIDs {
			tagID := tagIDs[j]
			if tagID == "" {
				continue
			}
			if score, ok := tagMap[tagID]; ok {
				scoreSum += Normalization(score, maxScore, minScore)
				cnt++
			}
		}
		if cnt > 0 {
			ret[item.Music.ID] = scoreSum / float64(cnt)
		}
	}
	return ret, nil
}

func getTagMaxAndMin(request *model.FeedRequest) (max, min float64) {
	return request.UserProfile.TagScore[request.UserProfile.TagIDs[0]], request.UserProfile.TagScore[request.UserProfile.TagIDs[len(request.UserProfile.TagIDs)-1]]
}
