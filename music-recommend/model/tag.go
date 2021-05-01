package model

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

type Tag struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Status     int32  `json:"status" db:"status"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	tableTag     = "tag"
	TagDelimiter = "+"
)

type tagMysql struct {
}

var onceTagMysql = tagMysql{}

func NewTagMysql() *tagMysql {
	return &onceTagMysql
}

func (t *tagMysql) Insert(tag *Tag) error {

	return t.doInsert(client(dbMusicRecommendNameTest), tag)
}

func (t *tagMysql) doInsert(db sqlx.Execer, tag *Tag) error {
	if tag == nil || tag.Name == "" {
		return nil
	}
	sql := fmt.Sprintf("insert ignore into %s(name,status,update_time) values(?,?,?)", tableTag)
	_, err := db.Exec(sql, tag.Name, tag.Status, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (t *tagMysql) SelectAll(ctx context.Context) (tags []Tag, err error) {
	sql := fmt.Sprintf("select * from %s", tableTag)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &tags, sql)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return
}

func (t *tagMysql) SelectTag(id int) (*Tag, error) {
	var ret Tag
	sql := fmt.Sprintf("select * from %s where id = ?", tableTag)
	err := client(dbMusicRecommendNameTest).Get(&ret, sql, id)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &ret, nil
}

func (t *tagMysql) SelectTagForName(name string) (*Tag, error) {
	var ret Tag
	sql := fmt.Sprintf("select * from %s where name = ?", tableTag)
	err := client(dbMusicRecommendNameTest).Get(&ret, sql, name)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &ret, nil
}
