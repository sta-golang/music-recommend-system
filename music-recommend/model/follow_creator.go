package model

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

type FollowCreator struct {
	ID         int    `json:"id" db:"id"`
	Status     int32  `json:"status" db:"status"`
	CreatorID  int    `json:"creator_id" db:"creator_id"`
	Username   int    `json:"username" db:"username"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	tableFollowCreator = "follow_creator"
)

type followCreatorMysql struct {
}

var onceFollowCreator = followCreatorMysql{}

func NewFollowCreatorMysql() *followCreatorMysql {
	return &onceFollowCreator
}

func (fc *followCreatorMysql) Insert(follow *FollowCreator) error {
	sql := fmt.Sprintf("insert into %s(creator_id,username,create_time,update_time) values(?,?,?,?)", tableFollowCreator)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, follow.CreatorID, follow.Username, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (fc *followCreatorMysql) Delete(username string, creatorID int) error {
	sql := fmt.Sprintf("Delete from %s where creator_id = ? and username = ?", tableFollowCreator)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, creatorID, username)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (fc *followCreatorMysql) SelectFollows(username string, pos, limit int) (ids []string, err error) {
	sql := fmt.Sprintf("select creator_id from %s where username = ? limit ?,?", tableFollowCreator)
	err = client(dbMusicRecommendNameTest).Select(&ids, sql, username, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (fc *followCreatorMysql) SelectFollowsOrderByCreateTime(username string, pos, limit int) (ids []string, err error) {
	sql := fmt.Sprintf("select creator_id from %s where username = ? ordey by create_time desc limit ?,?", tableFollowCreator)
	err = client(dbMusicRecommendNameTest).Select(&ids, sql, username, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}
