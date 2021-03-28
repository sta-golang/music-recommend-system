package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"strconv"
	"strings"
)

type Music struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Status      int32   `json:"status" db:"status"`
	Title       string  `json:"title" db:"title"`
	HotScore    float64 `json:"hot_score" db:"hot_score"`
	CreatorIDs  string  `json:"creator_ids" db:"creator_ids"`
	CreatorNames string `json:"creator_names" db:"creator_names"`
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
	MusicCreatorNameDelimiter = "/"
	insertSQLFmt = "insert ignore into %s values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
)

type musicMysql struct {
}

var onceMusicMysql = musicMysql{}

func NewMusicMysql() *musicMysql {
	return &onceMusicMysql
}

func (md *musicMysql) InsertMusic(music *Music) error {

	return md.doInsertMusic(client(dbMusicRecommendNameTest), music)
}

func (md *musicMysql) doInsertMusic(db sqlx.Execer, music *Music) error {

	sql := fmt.Sprintf(insertSQLFmt, tableMusic)
	_, err := db.Exec(sql, music.ID, music.Name, music.Status,
		music.Title, music.HotScore, music.CreatorIDs,music.CreatorNames, music.MusicUrl, music.PlayTime,
		music.TagIDs, music.TagNames, music.ImageUrl, music.PublishTime, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (md *musicMysql) InsertMusicAndCreatorMusic(music *Music) error {
	err := newMysqlTransaction().Transaction(func(tx *sqlx.Tx) error {
		err := NewMusicMysql().doInsertMusic(tx, music)
		if err != nil {
			log.Error(err)
			return err
		}
		split := strings.Split(music.CreatorIDs, CreatorDelimiter)
		if len(split) < 0 {
			return nil
		}
		for _, str := range split {
			if len(str) <= 0 {
				continue
			}
			creatorID, err := strconv.Atoi(str)
			if err != nil {
				log.Error(err)
				return err
			}
			err = NewCreatorToMusicMysql().doInsert(tx, &CreatorToMusic{
				CreatorID:  creatorID,
				MusicID:    music.ID,
				Status:     0,
				UpdateTime: "",
			})
			if err != nil {
				log.Error(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return err
}

func (md *musicMysql) UpdateMusic(music *Music) (affected bool, err error) {
	return md.doUpdateMusic(client(dbMusicRecommendNameTest), music)
}

func (md *musicMysql) UpdateMusicHotSource(id int, incr int32) error {
	sql := fmt.Sprintf("update %s set hot_score = hot_score + ? where id = ?", tableMusic)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, incr, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (md *musicMysql) SelectMusicForCreator(creatorID, pos, limit int) (musics []Music, err error) {
	sql := fmt.Sprintf("select * from %s where id in "+
		"(select music_id from %s where creator_id = ?) limit ?,?", tableMusic, tableCreatorToMusic)
	err = client(dbMusicRecommendNameTest).Select(&musics, sql, creatorID, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (md *musicMysql) SelectMusics(pos, limit int) (musics []Music, err error) {
	sql := fmt.Sprintf("select * from %s limit ?,?", tableMusic)
	err = client(dbMusicRecommendNameTest).Select(&musics, sql, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (md *musicMysql) FixMusicCreator(music *Music) error {
	if music == nil {
		return nil
	}
	sql := fmt.Sprintf("update %s set creator_ids=?, creator_names=?, update_time=? where id=?", tableMusic)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, music.CreatorIDs,music.CreatorNames, tm.GetNowDateTimeStr(), music.ID)
	if err != nil {
		log.Error(err)
		return  err
	}
	return nil
}

func (md *musicMysql) doUpdateMusic(db sqlx.Execer, music *Music) (affected bool, err error) {
	if music == nil {
		return false, nil
	}
	sql := fmt.Sprintf("update %s set status=?, hot_score=?, "+
		"tag_ids=?, tag_names=?, image_url=?,music_url=?,  update_time=? where id=?", tableMusic)
	res, err := db.Exec(sql, music.Status, music.HotScore, music.TagIDs,
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
