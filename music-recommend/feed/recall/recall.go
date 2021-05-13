package recall

import (
	"context"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/feed/recall/plugin"
	"github.com/sta-golang/music-recommend/model"
)

const (
	TagRecallPluginName = "tag"

	CreatorRecallPluginName        = "creator"
	ManMadePluginName              = "manMade"
	CreatorSimilarRecallPluginName = "creatorSimilar"
	maxRecallNum                   = 5000
)

type Handler func(*model.FeedRequest, string) ([]model.Music, error)

var defaultPlugins = map[string]model.PluginParams{
	TagRecallPluginName:            {PluginName: TagRecallPluginName, PluginParam: ""},
	CreatorRecallPluginName:        {PluginName: CreatorRecallPluginName, PluginParam: ""},
	CreatorSimilarRecallPluginName: {PluginName: CreatorSimilarRecallPluginName, PluginParam: ""},
	ManMadePluginName:              {PluginName: ManMadePluginName, PluginParam: ""},
}

var pluginMap = map[string]Handler{
	TagRecallPluginName:            plugin.TagRecall,
	CreatorRecallPluginName:        plugin.CreatorRecall,
	CreatorSimilarRecallPluginName: plugin.CreatorSimilarRecall,
	ManMadePluginName:              plugin.ManMadeRecall,
}

func DefaultRecallPlugin() map[string]model.PluginParams {
	return defaultPlugins
}

func FeedRecall(request *model.FeedRequest) error {
	log.InfoContext(request.Ctx, "Recall")
	if request.UserProfile == nil {
		return nil
	}
	if request.RecallPlugins == nil {
		request.RecallPlugins = DefaultRecallPlugin()
	}
	graph := buildDag(request)
	graph.Do(request.Ctx, false)
	items, err := graph.GetRootTask().GetRet()
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	if items == nil {
		return nil
	}
	request.RecallResults = items.([]model.Item)
	graph.DestoryAsync()
	return nil
}

func buildDag(request *model.FeedRequest) *dag.DagTasks {
	merge := dag.NewTask("merge", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		// TODO: 加入到sync.pool中  <06-05-21, FOUR SEASONS> //
		var ret []model.Item
		musicSet := set.NewHashSet(10000)
		for i := 0; i < helper.GetSubTaskSize(); i++ {
			if len(ret) >= maxRecallNum {
				break
			}
			subRet, err := helper.GetSubTaskRetForIndex(i)
			if err != nil {
				log.ErrorContext(request.Ctx, err)
				continue
			}
			if subRet == nil {
				continue
			}
			subMusics := subRet.([]model.Music)
			for j := 0; j < len(subMusics); j++ {
				if len(ret) >= maxRecallNum {
					break
				}
				if musicSet.Contains(subMusics[j].ID) {
					continue
				}
				musicSet.Add(subMusics[j].ID)
				ret = append(ret, model.Item{
					Music: subMusics[j],
					Score: 0.0,
				})
			}
		}
		log.InfoContextf(ctx, "Recall length : %d", len(ret))
		return ret, nil
	})
	for pluginName, recallPlugin := range request.RecallPlugins {
		if handler, ok := pluginMap[pluginName]; ok {
			curPluginName := pluginName
			merge.AddSubTask(dag.NewTask(pluginName, func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
				musics, err := handler(request, recallPlugin.PluginParam)
				if err != nil {
					log.ErrorContextf(ctx, "[cmd = %s] PluginName : %s err : %v", "Recall", pluginName, err)
					return nil, err
				}
				log.InfoContextf(request.Ctx, "Recall Plugin : %s racall length : %d", curPluginName, len(musics))
				return musics, err
			}))
		} else {
			log.ErrorContextf(request.Ctx, "[cmd = %s] PluginName : %s not found", "Profile", pluginName)
		}

	}
	return dag.NewDag(merge)
}
