package model

import (
	"fmt"
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
	sql := fmt.Sprintf("insert ignore into %s(name,status,update_time) values(?,?,?)", tableTag)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, tag.Name, tag.Status, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (t *tagMysql) SelectTag(id int) (ret *Tag, err error) {
	sql := fmt.Sprintf("select * from %s where id = ?", tableTag)
	err = client(dbMusicRecommendNameTest).Get(&ret, sql, id)
	return
}
