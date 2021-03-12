package data_load

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
	"strings"
	"sync"
)

const (
	setSize = 100000
)

var onceFunc sync.Once

type MysqlDataWriter struct {
	initFlag bool
	idSet    *set.StringSet
	tagMap   map[int]*set.StringSet
	tagSet   *set.StringSet
}

func (ml *MysqlDataWriter) Init() {
	onceFunc.Do(func() {
		if ml.tagSet == nil {
			ml.tagSet = set.NewStringSet(32)
		}
		if ml.idSet == nil {
			ml.idSet = set.NewStringSet(setSize)
		}
		if ml.tagMap == nil {
			ml.tagMap = make(map[int]*set.StringSet, setSize)
		}
		ml.initFlag = true
	})
}

func (ml *MysqlDataWriter) WriterMusic(musics []model.Music) error {
	if !ml.initFlag {
		ml.Init()
	}
	for _, music := range musics {
		err := ml.WriterTag(music.TagNames)
		if err != nil {
			log.Error(err)
			return err
		}
		split := strings.Split(music.TagIDs, model.TagDelimiter)
		if ml.idSet.Contains(fmt.Sprintf("%d", music.ID)) {

			if tagSet, ok := ml.tagMap[music.ID]; ok {
				oldLen := tagSet.Size()
				tagSet.Add(split...)
				if tagSet.Size() == oldLen {
					continue
				}
				affected, err := model.NewMusicMysql().UpdateMusic(&music)
				if err != nil {
					log.Error(err)
				}
				if !affected {
					log.Warnf("music id: %d update not affected", music.ID)
				}
			}
			continue
		}
		err = model.NewMusicMysql().InsertMusic(&music)
		if err != nil {
			log.Error(err)
		}
		ml.idSet.Add(fmt.Sprintf("%d", music.ID))
		ml.tagMap[music.ID] = set.NewStringSet(7)
		ml.tagMap[music.ID].Add(split...)
	}
	return nil
}

func (ml *MysqlDataWriter) WriterTag(tagNames string) error {
	if !ml.initFlag {
		ml.Init()
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
