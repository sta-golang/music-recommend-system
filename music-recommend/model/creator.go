package model

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

const (
	UnknownType   = 0
	SuperstarType = 9
	MusicianType  = 1
)

type Creator struct {
	ID             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Status         int32  `json:"status" db:"status"`
	ImageUrl       string `json:"image_url" db:"image_url"`
	Description    string `json:"description" db:"description"`
	SimilarCreator string `json:"similar_creator" db:"similar_creator"`
	FansNum        int    `json:"fans_num" db:"fans_num"`
	Type           int    `json:"type" db:"type"`
	UpdateTime     string `json:"update_time" db:"update_time"`
}

const (
	tableCreator = "creator"
)

type creatorMysql struct {
}

var onceCreatorMysql = creatorMysql{}

func NewCreatorMysql() *creatorMysql {
	return &onceCreatorMysql
}

func (cm *creatorMysql) Insert(c *Creator) error {
	if c == nil {
		return nil
	}
	sql := fmt.Sprintf("insert ignore into %s(id,name,status,image_url,description,"+
		"similar_creator,type,fans_num ,update_time) values(?,?,?,?,?,?,?,?,?)", tableCreator)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, c.ID, c.Name, c.Status, c.ImageUrl, c.Description,
		c.SimilarCreator, c.Type, c.FansNum, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (cm *creatorMysql) Update(c *Creator) (bool, error) {
	if c == nil {
		return false, nil
	}
	sql := fmt.Sprintf("update %s set name=?, status=?, image_url=?,"+
		" description=?,similar_creator=?,type=?,fans_num=?,update_time=? where id = ?", tableCreator)
	res, err := client(dbMusicRecommendNameTest).Exec(sql, c.Name, c.Status, c.ImageUrl, c.Description,
		c.SimilarCreator, c.Type, c.FansNum, tm.GetNowDateTimeStr(), c.ID)
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
