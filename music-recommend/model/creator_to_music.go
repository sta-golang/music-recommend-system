package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

type CreatorToMusic struct {
	ID         int    `json:"id" db:"id"`
	CreatorID  int    `json:"creator_id" db:"creator_id"`
	MusicID    int    `json:"music_id" db:"music_id"`
	Status     int32  `json:"status" db:"status"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	tableCreatorToMusic = "creator_music"
)

type creatorToMusicMysql struct {
}

var onceCreatorToMusicMysql = creatorToMusicMysql{}

func NewCreatorToMusicMysql() *creatorToMusicMysql {
	return &onceCreatorToMusicMysql
}

func (cm *creatorToMusicMysql) Insert(ctm *CreatorToMusic) error {
	return cm.doInsert(client(dbMusicRecommendNameTest), ctm)
}

func (cm *creatorToMusicMysql) doInsert(db sqlx.Execer, ctm *CreatorToMusic) error {
	sql := fmt.Sprintf("insert ignore into %s(creator_id,music_id,status,update_time) values(?,?,?,?)", tableCreatorToMusic)
	_, err := db.Exec(sql, ctm.CreatorID, ctm.MusicID, ctm.Status, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}
