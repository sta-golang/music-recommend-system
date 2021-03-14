package data_load

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/async/group"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/sta-golang/music-recommend/model"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MysqlDataWriter struct {
	tagSet *set.StringSet
}

func (ml *MysqlDataWriter) WriterMusic(musics []model.Music) error {

	for _, music := range musics {

		err := ml.WriterTag(music.TagNames)
		if err != nil {
			log.Error(err)
			return err
		}

		tagIDs, err := ml.getTagIDs(music.TagNames)
		if err != nil {
			log.Error(err)
			return err
		}

		queryObj, err := model.NewMusicMysql().SelectByID(music.ID)
		if err != nil {
			log.Error(err)
			return err
		}
		if queryObj == nil {
			music.TagIDs = strings.Join(tagIDs, model.TagDelimiter)
			err = model.NewMusicMysql().InsertMusic(&music)
			if err != nil {
				log.Error(err)
				return err
			}
			continue
		}

		stringSet := set.NewStringSet(10)
		stringSet.Add(strings.Split(queryObj.TagIDs, model.TagDelimiter)...)
		oldIDs := stringSet.Iterator()
		oldLen := stringSet.Size()
		stringSet.Add(tagIDs...)
		if stringSet.Size() == oldLen {
			continue
		}
		ids := stringSet.Iterator()
		names, err := ml.getTagNames(ids)
		if err != nil {
			log.Error(err)
			return err
		}
		music.TagNames = strings.Join(names, model.TagDelimiter)
		music.TagIDs = strings.Join(ids, model.TagDelimiter)
		log.Infof("music : %v oldIds %v, newIds : %v, newTagNames : %v", music.ID,
			strings.Join(oldIDs, model.TagDelimiter), music.TagIDs, music.TagNames)
		affected, err := model.NewMusicMysql().UpdateMusic(&music)
		if err != nil {
			log.Error(err)
		}
		if !affected {
			log.Warnf("music id: %d update not affected", music.ID)
		}

	}
	return nil
}

func (ml *MysqlDataWriter) WriterCreator(creators []model.Creator) error {

	for _, creator := range creators {
		err := model.NewCreatorMysql().Insert(&creator)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// LoadCreatorToMysql... 执行一次
func (ml *MysqlDataWriter) LoadCreatorToMysql(creators []model.Creator) error {
	asyncG := group.NewAsyncGroup(runtime.NumCPU())
	defer asyncG.Close()

	for _, creator := range creators {
		tempCreator := creator
		if addErr := asyncG.Add(fmt.Sprintf("%d", creator.ID), func() (interface{}, error) {
			curCreator := &tempCreator
			err := model.NewCreatorMysql().Insert(curCreator)
			if err != nil {
				log.Error("task ", err)
				return nil, err
			}

			return nil, nil
		}); addErr != nil {
			log.Error(addErr)
			return addErr
		}
	}
	asyncG.Wait()
	for _, tk := range asyncG.Iterator() {
		_, err := tk.Ret()
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (ml *MysqlDataWriter) LoadMusicToMysql(musics []APIMusicDetail) error {
	asyncG := group.NewAsyncGroup(runtime.NumCPU(), 100000)
	defer asyncG.Close()
	for _, detail := range musics {
		tempDetail := detail
		if addErr := asyncG.Add(fmt.Sprintf("%d", tempDetail.ID), func() (interface{}, error) {
			curDetail := &tempDetail
			creatorIDs := make([]string, 0, len(curDetail.AR))
			for _, ar := range curDetail.AR {
				creatorIDs = append(creatorIDs, fmt.Sprintf("%d", ar.CreatorID))
			}
			err := model.NewMusicMysql().InsertMusicAndCreatorMusic(&model.Music{
				ID:          curDetail.ID,
				Name:        curDetail.Name,
				Status:      0,
				Title:       curDetail.AL.TitleName,
				HotScore:    0,
				CreatorIDs:  strings.Join(creatorIDs, model.CreatorDelimiter),
				MusicUrl:    "",
				PlayTime:    curDetail.Dt,
				TagIDs:      "",
				TagNames:    "",
				ImageUrl:    curDetail.AL.TitleUrl,
				PublishTime: tm.ParseDataTimeToStr(time.Unix(int64(curDetail.PublishTime)/1000, 0)),
				UpdateTime:  "",
			})
			if err != nil {
				log.Error(err)
			}
			return nil, err
		}); addErr != nil {
			log.Error(addErr)
			return addErr
		}
	}
	asyncG.Wait()
	for _, tk := range asyncG.Iterator() {
		_, err := tk.Ret()
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (ml *MysqlDataWriter) getTagNames(ids []string) ([]string, error) {
	ret := make([]string, 0, len(ids))
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		tag, err := model.NewTagMysql().SelectTag(idInt)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		ret = append(ret, tag.Name)
	}
	return ret, nil
}

func (ml *MysqlDataWriter) getTagIDs(tagNames string) ([]string, error) {
	if len(tagNames) == 0 {
		return make([]string, 0), nil
	}
	split := strings.Split(tagNames, model.TagDelimiter)
	ret := make([]string, 0, len(split))

	for _, name := range split {
		tag, err := model.NewTagMysql().SelectTagForName(name)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		ret = append(ret, fmt.Sprintf("%d", tag.ID))
	}
	return ret, nil
}

func (ml *MysqlDataWriter) WriterTag(tagNames string) error {
	if len(tagNames) == 0 {
		return nil
	}
	if ml.tagSet == nil {
		ml.tagSet = set.NewStringSet(3000)
	}
	split := strings.Split(tagNames, model.TagDelimiter)
	for _, name := range split {
		if ml.tagSet.Contains(name) {
			continue
		}
		err := model.NewTagMysql().Insert(&model.Tag{
			Name:       name,
			Status:     0,
			UpdateTime: "",
		})
		if err != nil {
			log.Error(err)
			return err
		}
		ml.tagSet.Add(name)
	}
	return nil
}
