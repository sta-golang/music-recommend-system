package data_load

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/async/group"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/sta-golang/music-recommend/model"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

func (ml *MysqlDataWriter) FixMusic() error {
	offset := 0
	limit := 1000
	workerP := workerpool.NewWithQueueSize(runtime.NumCPU(), 50000)
	wg := sync.WaitGroup{}
	for {
		wg.Add(1)
		tempOffset := offset
		_ = workerP.Submit(func() {
			funCOffset := tempOffset
			defer wg.Done()
			musics, err := model.NewMusicMysql().SelectMusics(funCOffset, limit)
			if err != nil {
				log.Errorf("offset : tempOffset has err : %v",err)
			}
			for _, music := range musics {
				if len(music.CreatorIDs) <= 0 {
					continue
				}
				split := strings.Split(music.CreatorIDs, model.CreatorDelimiter)
				arr := make([]string, 0, len(split))
				for _, ss := range split {
					if ss == "" {
						continue
					}
					arr = append(arr, ss)
				}
				creators, err := model.NewCreatorMysql().SelectCreatorForIDs(arr)
				if err != nil {
					log.Errorf("music id : %d has err : %v",music.ID, err)
				}
				ids := make([]string, 0,len(creators))
				names := make([]string,0, len(creators))
				for _, creator := range creators {
					ids = append(ids, fmt.Sprintf("%d", creator.ID))
					names = append(names, creator.Name)
				}
				music.CreatorIDs = strings.Join(ids,model.CreatorDelimiter)
				music.CreatorNames = strings.Join(names, model.MusicCreatorNameDelimiter)
				err = model.NewMusicMysql().FixMusicCreator(&music)
				if err != nil {
					log.Errorf("music id : %d has err : %v",music.ID, err)
				}
			}
			log.ConsoleLogger.Infof("offset : %d finish", funCOffset)
		})
		//目前总条数 294138
		if offset > 294138 {
			break
		}
		offset += limit

	}
	wg.Wait()
	return nil
}

func (ml *MysqlDataWriter) FixCreator() error {
	details, err := model.NewCreatorMysql().SelectCreatorsDetails(0, 9999999)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, creator := range details {
		split := strings.Split(creator.SimilarCreator, model.CreatorDelimiter)
		newIDSet := set.NewStringSet()
		for _, creatorIDStr := range split {
			creatorID, _ := strconv.Atoi(creatorIDStr)
			selectCreator, err := model.NewCreatorMysql().SelectCreator(creatorID)
			if err != nil {
				log.Error(err)
				return err
			}
			if selectCreator == nil {
				continue
			}
			newIDSet.Add(creatorIDStr)

		}
		creator.SimilarCreator = strings.Join(newIDSet.Iterator(), model.CreatorDelimiter)
		_, err = model.NewCreatorMysql().Update(&creator)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// LoadCreatorToMysql... 执行一次
func (ml *MysqlDataWriter) LoadCreatorToMysql(creators []model.Creator) error {
	asyncG := group.NewAsyncGroup(runtime.NumCPU(), 8192)
	defer asyncG.Close()
	existSet := set.NewStringSet(len(creators) << 1)
	for _, creator := range creators {
		if existSet.Contains(fmt.Sprintf("%d", creator.ID)) {
			continue
		}
		existSet.Add(fmt.Sprintf("%d", creator.ID))
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

func (ml *MysqlDataWriter) LoadMusicToMysql(musics []APIMusicDetail)  error {
	asyncG := group.NewAsyncGroup(runtime.NumCPU(), 100000)
	defer asyncG.Close()

	for _, detail := range musics {
		tempDetail := detail
		if addErr := asyncG.Add(fmt.Sprintf("%d", tempDetail.ID), func() (interface{}, error) {
			curDetail := &tempDetail
			creatorIDs := set.NewStringSet(len(curDetail.AR))
			names := make([]string, 0)
			ids := make([]string, 0)
			for _, ar := range curDetail.AR {
				if creatorIDs.Contains(fmt.Sprintf("%d", ar.CreatorID)) || ar.CreatorID == 0 {
					continue
				}
				creatorIDs.Add(fmt.Sprintf("%d", ar.CreatorID))
				selectCreator, err := model.NewCreatorMysql().SelectCreator(ar.CreatorID)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				if selectCreator == nil {
					continue
				}
				names = append(names, ar.CreatorName)
				ids = append(ids, fmt.Sprintf("%d", ar.CreatorID))
			}
			err := model.NewMusicMysql().InsertMusicAndCreatorMusic(&model.Music{
				ID:          curDetail.ID,
				Name:        curDetail.Name,
				Status:      0,
				Title:       curDetail.AL.TitleName,
				HotScore:    0,
				CreatorIDs:  strings.Join(ids, model.CreatorDelimiter),
				CreatorNames: strings.Join(names, model.MusicCreatorNameDelimiter),
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

func (ml *MysqlDataWriter) LoadPlaylistForTag(res *APIPlaylistDetailResult) error {
	detail := res
	if len(detail.Playlist.Tags) <= 0 {
		return nil
	}
	tagNames := strings.Join(detail.Playlist.Tags, model.TagDelimiter)
	err := ml.WriterTag(tagNames)
	if err != nil {
		log.Error(err)
		return err
	}

	tagIDs, err := ml.getTagIDs(tagNames)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, song := range detail.Playlist.SongIDs {
		queryObj, err := model.NewMusicMysql().SelectByID(song.ID)
		if err != nil {
			log.Error(err)
			return err
		}
		if queryObj == nil {
			continue
		}
		stringSet := set.NewStringSet(10)
		if len(queryObj.TagIDs) != 0 {
			stringSet.Add(strings.Split(queryObj.TagIDs, model.TagDelimiter)...)
		}
		oldLen := stringSet.Size()
		if oldLen > 12 {
			continue
		}
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
		queryObj.TagNames = strings.Join(names, model.TagDelimiter)
		queryObj.TagIDs = strings.Join(ids, model.TagDelimiter)
		affected, err := model.NewMusicMysql().UpdateMusic(queryObj)
		if err != nil {
			log.Error(err)
		}
		if !affected {
			log.Warnf("music id: %d update not affected", queryObj.ID)
		}
	}
	return nil
}

func (ml *MysqlDataWriter) getTagNames(ids []string) ([]string, error) {
	ret := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
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
		if name == "" {
			continue
		}
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
		if name == "" {
			continue
		}
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
