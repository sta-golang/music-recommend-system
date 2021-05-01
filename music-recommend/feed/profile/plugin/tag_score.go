package plugin

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
)

func TagScore(ctx context.Context, dbProfile *model.DBProfile, profile *model.Profile, params string) error {
	if dbProfile == nil || profile == nil {
		return nil
	}
	if profile.TagScore == nil {
		profile.TagScore = make(map[string]float64)
	}
	tagIDS := strings.Split(dbProfile.TagScore, model.ProfileDelimiter)
	var err error
	for i := range tagIDS {
		index := strings.Index(tagIDS[i], model.ProfileSourceDelimiter)
		if index == -1 {
			profile.TagScore[tagIDS[i]] = 0.0
			profile.TagIDs = append(profile.TagIDs, tagIDS[i])
			log.ErrorContextf(ctx, "%s not delimiter", tagIDS[i])
			continue
		}
		tagID := tagIDS[i][:index]
		scoreStr := tagIDS[i][index+1:]
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			score = 0.0
			log.ErrorContext(ctx, err)
		}
		profile.TagScore[tagID] = score
		profile.TagIDs = append(profile.TagIDs, tagID)
	}
	sort.Slice(tagIDS, func(i, j int) bool {
		return profile.TagScore[tagIDS[i]] > profile.TagScore[tagIDS[j]]
	})
	return err
}
