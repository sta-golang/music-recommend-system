// Package model provides ...
package model

import (
	"context"
	"fmt"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/time"
)

type DBUserRead struct {
	ID         int    `json:"id" db:"id"`
	Status     int32  `json:"status" db:"status"`
	Username   string `json:"username" db:"username"`
	MusicRead  string `json:"music_read" db:"music_read"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	MusicReadDelimter = "|"
	tableUserRead     = "user_read"
)

type userReadMysql struct {
}

var onceUserReadMysql = userReadMysql{}

func NewUserReadMysql() *userReadMysql {
	return &onceUserReadMysql
}

func (um *userReadMysql) Insert(ctx context.Context, userRead *DBUserRead) error {
	sql := fmt.Sprintf("insert into %s(status, username, music_read,create_time, update_time) values(?,?,?,?,?)", tableUserRead)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, userRead.Status, userRead.Username,
		userRead.MusicRead, time.GetNowDateTimeStr(), time.GetNowDateTimeStr())
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return err
}

func (um *userReadMysql) Update(ctx context.Context, userRead *DBUserRead) error {
	sql := fmt.Sprintf("update %s set status = ?, music_read = ?, update_time = ? where username = ?", tableUserRead)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, userRead.Status, userRead.MusicRead,
		time.GetNowDateTimeStr(), userRead.Username)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return err
}

func (um *userReadMysql) SelectByUsername(ctx context.Context, username string) (*DBUserRead, error) {
	var ret DBUserRead
	sql := fmt.Sprintf("select * from %s where username = ?", tableUserRead)
	err := client(dbMusicRecommendNameTest).GetContext(ctx, &ret, sql, username)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return &ret, err
}
