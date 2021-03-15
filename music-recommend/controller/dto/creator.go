package dto

import "github.com/sta-golang/music-recommend/model"

/**
作用与视图层返回数据
 */
type CreatorAndSimilar struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Status         int32           `json:"status"`
	ImageUrl       string          `json:"image_url"`
	Description    string          `json:"description"`
	SimilarCreator []model.Creator `json:"similar_creator"`
	FansNum        int             `json:"fans_num"`
	Type           int             `json:"type"`
	UpdateTime     string          `json:"update_time"`
}

func NewCreatorAndSimilar(creator *model.Creator, similar []model.Creator) *CreatorAndSimilar {
	return &CreatorAndSimilar{
		ID:             creator.ID,
		Name:           creator.Name,
		Status:         creator.Status,
		ImageUrl:       creator.ImageUrl,
		Description:    creator.Description,
		SimilarCreator: similar,
		FansNum:        creator.FansNum,
		Type:           creator.Type,
		UpdateTime:     creator.UpdateTime,
	}
}
