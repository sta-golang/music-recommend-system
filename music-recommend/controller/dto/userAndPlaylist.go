package dto

import "github.com/sta-golang/music-recommend/model"

type UserAndPlaylist struct {
	Playlist *model.Playlist `json:"playlist"`
	User     *model.User     `json:"user"`
}
