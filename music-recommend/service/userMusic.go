package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/feed/profile"
	"github.com/sta-golang/music-recommend/model"
)

type userMusic struct {
	table map[string][]int
}

var PubUserMusicService = &userMusic{
	table: make(map[string][]int),
}

func (us *userMusic) RegisterStatistics() {
	PubStatisticsService.Register("userMusic", &StatisticsFunc{
		ParseFunc:   us.parseStatistics,
		ProcessFunc: us.processStatistics,
	})
}

func (us *userMusic) StatMusicForUser(ctx context.Context, user *model.User, musicID int) {
	if user == nil {
		return
	}
	if user.Username == "" {
		return
	}
	userStr := fmt.Sprintf("%s-%d", user.Username, musicID)
	us.parseStatistics(str.StringToBytes(&userStr))
}

func (us *userMusic) processStatistics() {
	ctx := context.Background()
	for key, val := range us.table {
		us.table[key] = nil
		insertFlag := false
		pf, err := profile.GetDBUserProfileWithCache(ctx, key)
		if err != nil {
			log.ErrorContextf(ctx, "username: %s val : %v update error : %v", key, val, err)
			continue
		}
		if pf == nil {
			insertFlag = true
			pf = &model.DBProfile{
				Username: key,
			}
			profile.SetProfileWithCache(pf)
		}
		err = us.processUserProfile(ctx, pf, val)
		if err != nil {
			log.ErrorContext(ctx, err)
			continue
		}
		if insertFlag {

			err = model.NewPeofileMysql().Insert(ctx, pf)
			if err != nil {
				log.ErrorContext(ctx, err)
				continue
			}
			log.Infof("username : %s profile insert", pf.Username)
		} else {
			err = model.NewPeofileMysql().Update(ctx, pf)
			if err != nil {
				log.ErrorContext(ctx, err)
				continue
			}
			log.Infof("username : %s profile update", pf.Username)
		}
	}
}

func (us *userMusic) processUserProfile(ctx context.Context, profile *model.DBProfile, musicIDS []int) error {
	var tagMap map[string]int
	if profile.TagScore == "" {
		tagMap = make(map[string]int)
	} else {
		err := codec.API.JsonAPI.Unmarshal(str.StringToBytes(&profile.TagScore), &tagMap)
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
	}
	for _, musicID := range musicIDS {
		music, sErr := PubMusicService.GetMusic(musicID)
		if sErr != nil {
			log.ErrorContext(ctx, sErr)
			return sErr.Err
		}
		if music == nil {
			continue
		}
		profile.MusicClick = fmt.Sprintf("%d%s%s", music.ID, model.ProfileDelimiter, profile.MusicClick)
		if music.TagIDs == "" {
			continue
		}
		tagIDS := strings.Split(music.TagIDs, model.TagDelimiter)
		for i := range tagIDS {
			cnt := 1
			if val, ok := tagMap[tagIDS[i]]; ok {
				cnt += val
			}
			tagMap[tagIDS[i]] = cnt
		}
	}
	bys, err := codec.API.JsonAPI.Marshal(tagMap)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	profile.TagScore = str.BytesToString(bys)

	return nil
}

func (us *userMusic) parseStatistics(bytes []byte) {
	userStr := str.BytesToString(bytes)
	log.Debug(userStr, "stat")
	split := strings.Split(userStr, "-")
	username := split[0]
	musicID, _ := strconv.Atoi(split[1])
	us.table[username] = append(us.table[username], musicID)
}
