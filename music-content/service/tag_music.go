package service

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
)

type tagMusicService struct {
}

var PubTagMusicService = &tagMusicService{}

func (tm *tagMusicService) ReLoadTagMusic(ctx context.Context) error {
	tags, err := model.NewTagMysql().SelectAll(ctx)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	tagMap := make(map[string][]model.Music)
	for _, tag := range tags {
		tagMap[strconv.Itoa(tag.ID)] = make([]model.Music, 0, 100000)
	}
	err = tm.loadForDB(ctx, tagMap)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	for _, musics := range tagMap {
		sort.Slice(musics, func(i, j int) bool {
			if musics[i].HotScore > musics[j].HotScore {
				return true
			}
			return false
		})
	}
	err = tm.saveTagMusicsDB(ctx, tagMap)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (tm *tagMusicService) saveTagMusicsDB(ctx context.Context, tagMap map[string][]model.Music) error {
	for key, musics := range tagMap {

		tagID, err := strconv.Atoi(key)
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		tgs := &model.TagMusic{
			TagID: tagID,
		}
		strArr := make([]string, 0, len(musics))
		for _, music := range musics[:minInt(5000, len(musics))] {
			strArr = append(strArr, strconv.Itoa(music.ID))
		}
		tgs.Musics = strings.Join(strArr, model.TagMusicDelimter)
		affected, err := model.NewTagMusicMysql().Insert(ctx, tgs)
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		if !affected {
			err = model.NewTagMusicMysql().Update(ctx, tgs)
			if err != nil {
				log.ErrorContext(ctx, err)
				return err
			}
		}
	}
	return nil
}

func (tm *tagMusicService) loadForDB(ctx context.Context, tagMap map[string][]model.Music) error {
	if tagMap == nil {
		return nil
	}
	start, limit := 0, 3000
	for {
		musics, err := model.NewMusicMysql().SelectMusicsWithContent(start, limit)
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		if len(musics) <= 0 {
			break
		}
		for _, music := range musics {
			if music.TagIDs == "" {
				continue
			}
			tagids := strings.Split(music.TagIDs, model.MusicTagDelimiter)
			for _, tagID := range tagids {
				tagMap[tagID] = append(tagMap[tagID], music)
			}
		}
		start += limit
	}

	return nil
}
