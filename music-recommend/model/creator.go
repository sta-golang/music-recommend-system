package model

const (
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
	Type           int    `json:"type" db:"type"`
	UpdateTime     string `json:"update_time" db:"update_time"`
}
