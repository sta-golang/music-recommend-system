package rank

import (
	"context"
	"sort"

	"github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/feed/rank/plugin"
	"github.com/sta-golang/music-recommend/model"
)

const (
	TagScore     = "tagScore"
	CreatorScore = "creatorScore"
	ManMadeScore = "manMadeScore"
)

var defautlPlugin = map[string]model.PluginParams{
	TagScore:     {PluginName: TagScore, PluginParam: ""},
	CreatorScore: {PluginName: CreatorScore, PluginParam: ""},
	ManMadeScore: {PluginName: ManMadeScore, PluginParam: ""},
}

func DefaultRankPlugin() map[string]model.PluginParams {
	return defautlPlugin
}

var pluginMap = map[string]Handler{
	TagScore:     plugin.TagRank,
	CreatorScore: plugin.CreatorRank,
	ManMadeScore: plugin.ManMadeRank,
}

type Handler func(*model.FeedRequest, string) (map[int]float64, error)

func FeedRank(request *model.FeedRequest) error {
	log.InfoContext(request.Ctx, "Rank")
	if len(request.RecallResults) <= 0 {
		return nil
	}
	if request.RankPlugins == nil {
		request.RankPlugins = DefaultRankPlugin()
	}
	graph := buildDag(request)
	graph.Do(request.Ctx, false)
	request.RankResults = request.RecallResults
	request.RecallResults = nil
	graph.DestoryAsync()
	return nil
}

func buildDag(request *model.FeedRequest) *dag.DagTasks {
	merge := dag.NewTask("merge", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		scoreMap := make(map[int]float64)
		for i := 0; i < helper.GetSubTaskSize(); i++ {
			subRet, subErr := helper.GetSubTaskRetForIndex(i)
			if subErr != nil {
				log.ErrorContext(ctx, subErr)
			}
			if subRet == nil {
				continue
			}
			subScoreMap := subRet.(map[int]float64)
			for key, val := range subScoreMap {
				scoreMap[key] += val
			}
		}
		for i := range request.RecallResults {
			request.RecallResults[i].Score = scoreMap[request.RecallResults[i].Music.ID]
		}
		sort.Slice(request.RecallResults, func(i, j int) bool {
			if request.RecallResults[i].Score > request.RecallResults[j].Score {
				return true
			}
			return false
		})
		return nil, nil
	})
	for pluginName, recallPlugin := range request.RankPlugins {
		if handler, ok := pluginMap[pluginName]; ok {
			merge.AddSubTask(dag.NewTask(pluginName, func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
				musics, err := handler(request, recallPlugin.PluginParam)
				if err != nil {
					log.ErrorContextf(ctx, "[cmd = %s] PluginName : %s err : %v", "Rank", pluginName, err)
					return nil, err
				}
				return musics, err
			}))

		} else {
			log.ErrorContextf(request.Ctx, "[cmd = %s] PluginName : %s not found", "Profile", pluginName)
		}

	}
	return dag.NewDag(merge)
}
