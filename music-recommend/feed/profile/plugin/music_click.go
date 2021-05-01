package plugin

import (
	"context"
	"strings"

	"github.com/sta-golang/music-recommend/model"
)

func MusicClick(ctx context.Context, dbProfile *model.DBProfile, profile *model.Profile, params string) error {
	if dbProfile == nil {
		return nil
	}
	profile.MusicClick = strings.Split(dbProfile.MusicClick, model.ProfileDelimiter)
	return nil
}
