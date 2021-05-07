package model

import "context"

type FeedRequest struct {
	Ctx context.Context

	AnyUser            bool
	Username           string
	User               *User
	UserProfilePlugins map[string]PluginParams
	UserProfile        *Profile
	RecallPlugins      map[string]PluginParams
	RecallResults      []Item
	RankPlugins        map[string]PluginParams
	RankResults        []Item
	FeedResults        []Item
}

type PluginParams struct {
	PluginName  string `json:"pluginName"`
	PluginParam string `json:"pluginParam"`
}

type Item struct {
	Music Music   `json:"music"`
	Score float64 `json:"score"`
	Extra string  `json:"extra"`
}
