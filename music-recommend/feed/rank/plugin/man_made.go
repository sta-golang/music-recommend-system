// Package plugin provides ...
package plugin

import (
	"strconv"

	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
)

func ManMadeRank(request *model.FeedRequest, params string) (map[int]float64, error) {
	ret := make(map[int]float64)
	for i := range request.RecallResults {
		item := request.RecallResults[i]
		musicID := strconv.Itoa(item.Music.ID)
		if service.ContainsManMade(musicID) {
			ret[item.Music.ID] = 0.5
		}
	}
	return ret, nil
}
