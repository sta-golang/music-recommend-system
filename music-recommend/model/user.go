package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"time"
)

type User struct {
	ID int `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Name string `json:"name" db:"name"`
	Status int32 `json:"status" db:"status"`
	ImageUrl string `json:"image_url" db:"image_url"`
	CreateTime string `json:"create_time" db:"create_time"`
	LastLoginTime string `json:"last_login_time" db:"last_login_time"`
	LastMonthLoginNum int `json:"last_month_login_num" db:"last_month_login_num"`
	LastStatTime string `json:"last_stat_time" db:"last_stat_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

const (
	tableUser = "user"
)
type userMysql struct {
}

var onceUserMysql = userMysql{}

func NewUserMysql() *userMysql {
	return &onceUserMysql
}

func (um *userMysql) Insert(user *User) error {
	return um.doInsert(client(dbMusicRecommendNameTest), user)
}

func (um *userMysql) doInsert(db sqlx.Execer, user *User) error {
	sql := fmt.Sprintf("insert into %s(username,password,name,status,image_url,create_time,last_login_time," +
		"last_month_login_num,last_stat_time,update_time) values(?,?,?,?,?,?,?,?,?,?)", tableUser)
	_, err := db.Exec(sql, user.Username, user.Password, user.Name, 0, user.ImageUrl, tm.GetNowDateTimeStr(),
		tm.GetNowDateTimeStr(), 0, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (um *userMysql) Update(user *User) (bool, error) {
	return um.doUpdate(client(dbMusicRecommendNameTest), user)
}

func (um *userMysql) doUpdate(db sqlx.Execer, user *User) (bool, error) {
	sql := fmt.Sprintf("update %s set name=?,image_url=?,status=?,update_time=? where id = ?", tableUser)
	res, err := db.Exec(sql, user.Name,user.ImageUrl,user.Status,tm.GetNowDateTimeStr(), user.ID)
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

func (um *userMysql) StatLogin(user *User) error {
	sql := fmt.Sprintf("update %s set last_login_time=?,last_month_login_num=?last_stat_time=?,update_time", tableUser)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, user.LastLoginTime, user.LastMonthLoginNum, user.LastStatTime, tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (um *userMysql) ReSetStatistics(username string) (bool, error) {
	sql := fmt.Sprintf("update %s set last_month_login_num=?,last_stat_time=?,last_login_time=?," +
		"update_time=? where username=? and last_stat_time <= ?", tableUser)
	res, err := client(dbMusicRecommendNameTest).Exec(sql, 1, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr(),
		tm.GetNowDateTimeStr(), username, tm.ParseDataTimeToStr(tm.GetNowTime().Add(-(time.Hour * 24 * 30))))
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

func (um *userMysql) UserLogin(username string) error {
	sql := fmt.Sprintf("update %s set last_month_login_num=last_month_login_num+1,last_login_time=?," +
		"update_time=? where username=? ", tableUser)
	_, err := client(dbMusicRecommendNameTest).Exec(sql, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr(), username)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (um *userMysql) SelectUser(id int) (*User, error) {
	var ret User
	sql := fmt.Sprintf("select * from %s where id = ?", tableUser)
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

func (um *userMysql) SelectUserForUserName(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	var ret User
	sql := fmt.Sprintf("select * from %s where username = ?", tableUser)
	err := client(dbMusicRecommendNameTest).Get(&ret, sql, username)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &ret, nil
}