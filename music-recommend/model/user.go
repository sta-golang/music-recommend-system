package model

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
