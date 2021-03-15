package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"strings"
)

const (
	CreatorUnknownType   = 0
	CreatorSuperstarType = 9
	CreatorMusicianType  = 1
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
	tableCreator     = "creator"
	CreatorDelimiter = "+"

	StatusLoadMusicFinish = 1
)

type creatorMysql struct {
}

var onceCreatorMysql = creatorMysql{}

func NewCreatorMysql() *creatorMysql {
	return &onceCreatorMysql
}

func (cm *creatorMysql) Insert(c *Creator) error {
	return cm.doInsert(client(dbMusicRecommendNameTest), c)
}

func (cm *creatorMysql) doInsert(db sqlx.Execer, c *Creator) error {
	if c == nil {
		return nil
	}
	sql := fmt.Sprintf("insert ignore into %s(id,name,status,image_url,description,"+
		"similar_creator,type,fans_num ,update_time) values(?,?,?,?,?,?,?,?,?)", tableCreator)
	_, err := db.Exec(sql, c.ID, c.Name, c.Status, c.ImageUrl, c.Description,
		c.SimilarCreator, c.Type, c.FansNum, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (cm *creatorMysql) SelectCreator(id int) (*Creator, error) {
	var ret Creator
	sql := fmt.Sprintf("select * from %s where id = ?", tableCreator)
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

func (cm *creatorMysql) SelectCreatorForIDs(ids []string) (creators []Creator, err error) {
	sql := fmt.Sprintf("select id, name, image_url,type from %s where id in(%s)",
		tableCreator, strings.Join(ids, ","))
	err = client(dbMusicRecommendNameTest).Select(&creators, sql)
	if err != nil {
		log.Error("ids : %v has err : %v", ids, err)
		return nil, err
	}
	return
}

func (cm *creatorMysql) SelectCreators(pos, limit int) (creators []Creator, err error) {
	sql := fmt.Sprintf("select id, name, image_url,type from %s limit ?,?", tableCreator)
	err = client(dbMusicRecommendNameTest).Select(&creators, sql, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (cm *creatorMysql) SelectCreatorsDetails(pos, limit int) (creators []Creator, err error) {
	sql := fmt.Sprintf("select * from %s limit ?,?", tableCreator)
	err = client(dbMusicRecommendNameTest).Select(&creators, sql, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (cm *creatorMysql) SelectCreatorsForStatus(status, pos, limit int) (creators []Creator, err error) {
	sql := fmt.Sprintf("select id, name, image_url,type from %s where status = ? limit ?,?", tableCreator)
	err = client(dbMusicRecommendNameTest).Select(&creators, sql,status, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (cm *creatorMysql) UpdateCreatorsForStatus(status int32, id int) (bool, error) {
	sql := fmt.Sprintf("update %s set status = ? where id = ?", tableCreator)
	res, err := client(dbMusicRecommendNameTest).Exec(sql, status, id)
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

func (cm *creatorMysql) SelectCreatorsForType(ty, pos, limit int) (creators []Creator, err error) {
	sql := fmt.Sprintf("select id, name, image_url,type from %s where type = ? limit ?,?", tableCreator)
	err = client(dbMusicRecommendNameTest).Select(&creators, sql, ty, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (cm *creatorMysql) Update(c *Creator) (bool, error) {
	return cm.doUpdate(client(dbMusicRecommendNameTest), c)
}
func (cm *creatorMysql) doUpdate(db sqlx.Execer, c *Creator) (bool, error) {
	if c == nil {
		return false, nil
	}
	sql := fmt.Sprintf("update %s set name=?, status=?, image_url=?,"+
		" description=?,similar_creator=?,type=?,fans_num=?,update_time=? where id = ?", tableCreator)
	res, err := db.Exec(sql, c.Name, c.Status, c.ImageUrl, c.Description,
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
