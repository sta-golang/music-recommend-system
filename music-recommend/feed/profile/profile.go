package profile

import (
	"context"
	"fmt"

	"github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/feed/profile/plugin"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
)

const (
	dbProfileCacheKey = "dbProfile_%s"
)

type ProfileParams struct {
	Username string
	Plugins  map[string]model.PluginParams
}

var defaultPlugins = map[string]model.PluginParams{
	"musicClick": {PluginName: "musicClick", PluginParam: ""},
	"tagScore":   {PluginName: "tagScore", PluginParam: ""},
}

func DefaultParams(username string) ProfileParams {
	return ProfileParams{
		Username: username,
		Plugins:  defaultPlugins,
	}
}

func NewParams(username string, plugins map[string]model.PluginParams) ProfileParams {
	return ProfileParams{
		Username: username,
		Plugins:  plugins,
	}
}

func GetDBUserProfile(ctx context.Context, username string) (*model.DBProfile, error) {
	dbProfile, err := model.NewPeofileMysql().SelectByUsername(ctx, username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return dbProfile, nil
}

func SetProfileWithCache(profile *model.DBProfile) {
	if profile == nil {
		return
	}
	key := fmt.Sprintf(dbProfileCacheKey, profile.Username)
	cache.PubCacheService.Set(key, profile, cache.Hour*12, cache.Nine)
}

func GetDBUserProfileWithCache(ctx context.Context, username string) (*model.DBProfile, error) {
	key := fmt.Sprintf(dbProfileCacheKey, username)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.(*model.DBProfile), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ret, err := GetDBUserProfile(ctx, username)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			cache.PubCacheService.Set(key, ret, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, ret, cache.Hour*12, cache.Nine)
		return ret, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*model.DBProfile), nil
}

type Handler func(request *model.FeedRequest, dbProfile *model.DBProfile, params string) error

var pluginMap = map[string]Handler{
	"musicClick": plugin.MusicClick,
	"tagScore":   plugin.TagScore,
}

func GetDefaultProfile() *model.Profile {
	return &model.Profile{
		TagIDs: []string{"1", "2", "27", "3", "7"},
		TagScore: map[string]float64{
			"1":  1.0,
			"2":  1.0,
			"27": 1.0,
			"3":  1.0,
			"7":  1.0,
		},
	}
}

func FeedUserProfile(request *model.FeedRequest) error {
	log.InfoContext(request.Ctx, "Profile")
	if request.Username == "" {
		return fmt.Errorf("用户出现异常")
	}
	if request.AnyUser {
		request.UserProfile = GetDefaultProfile()
		return nil
	}
	request.UserProfile = &model.Profile{}
	var dbProfile *model.DBProfile
	var err error
	if !request.AnyUser {
		dbProfile, err = GetDBUserProfileWithCache(request.Ctx, request.Username)
		if err != nil {
			log.ErrorContext(request.Ctx, err)
		}
	}
	graph := buildDag(request, request.UserProfilePlugins, dbProfile)
	graph.Do(request.Ctx, false)
	graph.DestoryAsync()
	return nil
}

func buildDag(request *model.FeedRequest, plugins map[string]model.PluginParams, dbProfile *model.DBProfile) *dag.DagTasks {
	rootTask := dag.NewTask("root", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		return nil, nil
	})
	for pluginName, pluginParam := range plugins {
		if handler, ok := pluginMap[pluginName]; ok {
			rootTask.AddSubTask(dag.NewTask(pluginName, func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
				err := handler(request, dbProfile, pluginParam.PluginParam)
				if err != nil {
					log.ErrorContextf(ctx, "[cmd = %s] PluginName : %s err : %v", "Profile", pluginName, err)
				}
				return nil, err
			}))

		} else {
			log.ErrorContextf(request.Ctx, "[cmd = %s] PluginName : %s not found", "Profile", pluginName)
		}
	}
	return dag.NewDag(rootTask)
}
