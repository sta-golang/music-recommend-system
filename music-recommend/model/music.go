package model

type Music struct {
	ID int `json:"id" db:"id"`
	RefID int `json:"ref_id" db:"ref_id"`
	Name string `json:"name" db:"name"`
	Status int32 `json:"status" db:"status"`
	Title string `json:"title" db:"title"`
	CreatorID int `json:"creator_id" db:"creator_id"`
	CreatorName string `json:"creator_name" db:"creator_name"`
	PlayTime int `json:"play_time" db:"play_time"`
	ImageUrl string `json:"image_url" db:"image_url"`
	PublishTime string `json:"publish_time" db:"publish_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}
