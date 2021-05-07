package model

import (
	"context"
	"fmt"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/time"
)

type TagMusic struct {
	ID         int    `json:"id" db:"id"`
	TagID      int    `json:"tag_id" db:"tag_id"`
	Musics     string `json:"musics" db:"musics"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	tableTagMusic    = "tag_music"
	TagMusicDelimter = "+"
)

type tagMusicMysql struct {
}

var onceTagMusicMysql = tagMusicMysql{}

func NewTagMusicMysql() *tagMusicMysql {
	return &onceTagMusicMysql
}

func (tm *tagMusicMysql) Select(ctx context.Context) (tgs []TagMusic, err error) {
	sql := fmt.Sprintf("select * from %s", tableTagMusic)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &tgs, sql)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return tgs, nil
}

func (tm *tagMusicMysql) SelectWithTagID(ctx context.Context, tagID int) (*TagMusic, error) {
	var ret TagMusic
	sql := fmt.Sprintf("select * from %s where tag_id = ?", tableTagMusic)
	err := client(dbMusicRecommendNameTest).GetContext(ctx, &ret, sql, tagID)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return &ret, nil
}

func (tm *tagMusicMysql) Insert(ctx context.Context, tgs *TagMusic) (bool, error) {
	if tgs == nil {
		return false, nil
	}
	sql := fmt.Sprintf("insert ignore into %s(tag_id, musics, create_time, update_time) values(?,?,?,?)", tableTagMusic)
	res, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, tgs.TagID, tgs.Musics, time.GetNowDateTimeStr(), time.GetNowDateTimeStr())
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}
	return rows > 0, err
}

func (tm *tagMusicMysql) Update(ctx context.Context, tgs *TagMusic) error {
	if tgs == nil {
		return nil
	}
	sql := fmt.Sprintf("update %s set musics = ?, update_time = ? where tag_id = ?", tableTagMusic)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, tgs.Musics, time.GetNowDateTimeStr(), tgs.TagID)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	return nil
}
