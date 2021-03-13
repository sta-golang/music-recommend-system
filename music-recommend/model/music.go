package model

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

type Music struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Status      int32   `json:"status" db:"status"`
	Title       string  `json:"title" db:"title"`
	HotScore    float64 `json:"hot_score" db:"hot_score"`
	CreatorID   int     `json:"creator_id" db:"creator_id"`
	MusicUrl    string  `json:"music_url" db:"music_url"`
	PlayTime    int     `json:"play_time" db:"play_time"`
	TagIDs      string  `json:"tag_ids" db:"tag_ids"`
	TagNames    string  `json:"tag_names" db:"tag_names"`
	ImageUrl    string  `json:"image_url" db:"image_url"`
	PublishTime string  `json:"publish_time" db:"publish_time"`
	UpdateTime  string  `json:"update_time" db:"update_time"`
}

const (
	tableMusic                = "music"
	MusicDefaultStatus        = 0
	MusicHasUrlMusicUrlStatus = 1
)

type musicMysql struct {
}

var onceMusicMysql = musicMysql{}

func NewMusicMysql() *musicMysql {
	return &onceMusicMysql
}

func (md *musicMysql) InsertMusic(music *Music) error {
	if music == nil {
		return nil
	}
	sql := fmt.Sprintf("insert ignore into %s values(?,?,?,?,?,?,?,?,?,?,?,?,?)", tableMusic)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, music.ID, music.Name, music.Status,
		music.Title, music.HotScore, music.CreatorID, music.ImageUrl, music.TagIDs, music.TagNames,
		music.PlayTime, music.ImageUrl, music.PublishTime, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (md *musicMysql) UpdateMusic(music *Music) (affected bool, err error) {
	if music == nil {
		return false, nil
	}
	sql := fmt.Sprintf("update %s set status=?, hot_score=?, "+
		"tag_ids=?, tag_names=?, image_url=?,music_url=?, update_time=? where id=?", tableMusic)
	res, err := client(dbMusicRecommendNameTest).Exec(sql, music.Status, music.HotScore, music.TagIDs,
		music.TagNames, music.ImageUrl, music.MusicUrl, tm.GetNowDateTimeStr(), music.ID)
	if err != nil {
		log.Error(err)
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return false, err
	}
	return rows > 0, err
}

func (md *musicMysql) SelectByID(id int) (*Music, error) {
	var ret Music
	sql := fmt.Sprintf("select * from %s where id = ?", tableMusic)
	err := client(dbMusicRecommendNameTest).Get(&ret, sql, id)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.Error(err)
	}
	return &ret, err
}
