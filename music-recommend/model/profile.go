package model

import (
	"context"
	"fmt"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/time"
)

type Profile struct {
	Username   string             `json:"username"`
	TagScore   map[string]float64 `json:"tag_score"`
	MusicClick []string           `json:"music_click"`
	TagIDs     []string           `json:"tag_ids"` // Tagid的倒排索引
}
type DBProfile struct {
	ID         int    `json:"id" db:"id"`
	Status     int32  `json:"status" db:"status"`
	Username   string `json:"username" db:"username"`
	MusicClick string `json:"music_click" db:"music_click"`
	TagScore   string `json:"tag_score" db:"tag_score"`
	Extra      string `json:"extra" db:"extra"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	tableProfile           = "profile"
	ProfileDelimiter       = "|"
	ProfileSourceDelimiter = ":"
)

type profileMysql struct {
}

var onceProfileMysql = profileMysql{}

func NewPeofileMysql() *profileMysql {
	return &onceProfileMysql
}

func (pm *profileMysql) Insert(ctx context.Context, profile *DBProfile) error {
	sql := fmt.Sprintf("insert into %s(status, username, music_click, tag_score, extra, create_time, update_time) values(?,?,?,?,?,?,?)", tableProfile)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, profile.Status, profile.Username,
		profile.MusicClick, profile.TagScore, profile.Extra, time.GetNowDateTimeStr(), time.GetNowDateTimeStr())
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return err
}

func (pm *profileMysql) Update(ctx context.Context, profile *DBProfile) error {
	sql := fmt.Sprintf("update %s set status = ?, music_click = ?, tag_score = ?, extra = ?, create_time = ?, update_time = ? where username = ?", tableProfile)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, profile.Status, profile.MusicClick, profile.TagScore, profile.Extra,
		time.GetNowDateTimeStr(), time.GetNowDateTimeStr(), profile.Username)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return err
}

func (pm *profileMysql) SelectByUsername(ctx context.Context, username string) (*DBProfile, error) {
	var ret DBProfile
	sql := fmt.Sprintf("select * from %s where username = ?", tableProfile)
	err := client(dbMusicRecommendNameTest).GetContext(ctx, &ret, sql, username)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return &ret, err
}
