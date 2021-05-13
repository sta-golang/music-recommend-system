package plugin

import (
	"sort"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/model"
)

func TagScore(request *model.FeedRequest, dbProfile *model.DBProfile, params string) error {
	if request == nil || dbProfile == nil {
		return nil
	}
	if request.UserProfile.TagScore == nil {
		request.UserProfile.TagScore = make(map[string]float64)
	}
	var err error
	tagSum := 0
	var tagMap map[string]int
	err = codec.API.JsonAPI.Unmarshal(str.StringToBytes(&dbProfile.TagScore), &tagMap)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	for key, val := range tagMap {
		tagSum += val
		request.UserProfile.TagIDs = append(request.UserProfile.TagIDs, key)
	}
	for key, val := range tagMap {
		request.UserProfile.TagScore[key] = float64(val) / float64(tagSum)
	}
	sort.Slice(request.UserProfile.TagIDs, func(i, j int) bool {
		return request.UserProfile.TagScore[request.UserProfile.TagIDs[i]] > request.UserProfile.TagScore[request.UserProfile.TagIDs[j]]
	})
	return nil

}
