// Package dto provides ...
package dto

import "github.com/sta-golang/music-recommend/model"

type RecommendResults struct {
	Results   []model.Item `json:"result"`
	SessionID string       `json:"sessionID"`
}
