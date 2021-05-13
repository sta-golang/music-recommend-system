package plugin

import (
	"strings"

	"github.com/sta-golang/music-recommend/model"
)

func MusicClick(request *model.FeedRequest, dbProfile *model.DBProfile, params string) error {
	if dbProfile == nil {
		return nil
	}
	request.UserProfile.MusicClick = strings.Split(dbProfile.MusicClick, model.ProfileDelimiter)
	return nil
}
