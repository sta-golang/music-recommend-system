package profile

import (
	"context"

	"github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/feed/profile/plugin"
	"github.com/sta-golang/music-recommend/model"
)

func init() {
	dag.Config().SetPool()
}

type ProfileParams struct {
	Username string
	Plugins  map[string]ProfilePlugin
}

type ProfilePlugin struct {
	PluginName   string
	PluginParams string
}

type Handler func(ctx context.Context, dbProfile *model.DBProfile, profile *model.Profile, params string) error

var pluginMap = map[string]Handler{
	"musicClick": plugin.MusicClick,
	"tagScore":   plugin.TagScore,
}

func GetUserProfile(ctx context.Context, params ProfileParams) (*model.Profile, error) {
	if params.Username == "" {
		return nil, nil
	}
	dbProfile, err := model.NewPeofileMysql().SelectByUsername(ctx, params.Username)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	if dbProfile == nil {
		return nil, nil
	}
	profile := &model.Profile{}
	graph := buildDag(ctx, &params, dbProfile, profile)
	graph.Do(ctx, false)
	graph.DestoryAsync()
	return nil, nil
}

func buildDag(ctx context.Context, params *ProfileParams, dbProfile *model.DBProfile, profile *model.Profile) *dag.DagTasks {
	rootTask := dag.NewTask("root", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		return nil, nil
	})
	for pluginName, pluginParam := range params.Plugins {
		if handler, ok := pluginMap[pluginName]; ok {
			rootTask.AddSubTask(dag.NewTask(pluginName, func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
				err := handler(ctx, dbProfile, profile, pluginParam.PluginParams)
				if err != nil {
					log.ErrorContextf(ctx, "[cmd = %s] PluginName : %s err : %v", "Profile", pluginName, err)
				}
				return nil, err
			}))

		} else {
			log.ErrorContextf(ctx, "[cmd = %s] PluginName : %s not found", "Profile", pluginName)
		}
	}
	return dag.NewDag(rootTask)
}
