package model

import (
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

type Music struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Status int32 `json:"status" db:"status"`
	Title string `json:"title" db:"title"`
	HotScore float64 `json:"hot_score" db:"hot_score"`
	CreatorID int `json:"creator_id" db:"creator_id"`
	CreatorName string `json:"creator_name" db:"creator_name"`
	PlayTime int `json:"play_time" db:"play_time"`
	ImageUrl string `json:"image_url" db:"image_url"`
	PublishTime string `json:"publish_time" db:"publish_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	dbMusicRecommendNameTest = "music_recommend_test"
)

type musicDB struct {
}

var onceMusicDB = musicDB{}

func NewMusicDB() *musicDB {
	return &onceMusicDB
}

func (md *musicDB) InsertMusic(music *Music) error {
	sql := "insert ignore into music values(?,?,?,?,?,?,?,?,?,?,?)"
	_, err := client(dbMusicRecommendNameTest).Exec(sql,music.ID, music.Name,music.Status,music.Title,music.HotScore,
		music.CreatorID,music.CreatorName,music.PlayTime, music.ImageUrl,music.PublishTime, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}