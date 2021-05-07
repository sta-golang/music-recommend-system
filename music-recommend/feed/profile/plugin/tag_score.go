package plugin

import (
	"context"
	"sort"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/model"
)

func TagScore(ctx context.Context, dbProfile *model.DBProfile, profile *model.Profile, params string) error {
	if profile == nil || dbProfile == nil {
		return nil
	}
	if profile.TagScore == nil {
		profile.TagScore = make(map[string]float64)
	}
	var err error
	tagSum := 0
	var tagMap map[string]int
	err = codec.API.JsonAPI.Unmarshal(str.StringToBytes(&dbProfile.TagScore), &tagMap)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	for key, val := range tagMap {
		tagSum += val
		profile.TagIDs = append(profile.TagIDs, key)
	}
	for key, val := range tagMap {
		profile.TagScore[key] = float64(val) / float64(tagSum)
	}
	sort.Slice(profile.TagIDs, func(i, j int) bool {
		return profile.TagScore[profile.TagIDs[i]] > profile.TagScore[profile.TagIDs[j]]
	})
	return nil

}
