package feed

import (
	"fmt"

	_ "github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/feed/profile"
	"github.com/sta-golang/music-recommend/feed/rank"
	"github.com/sta-golang/music-recommend/feed/recall"
	"github.com/sta-golang/music-recommend/feed/rerank"
	"github.com/sta-golang/music-recommend/model"
)

var (
	requestNilErr = fmt.Errorf("request not be nil")
)

func init() {
	//dag.Config().SetPool()
}

func FeedList(request *model.FeedRequest) error {
	if request == nil {
		return requestNilErr
	}
	if len(request.UserProfilePlugins) <= 0 {
		request.UserProfilePlugins = profile.DefaultParams(request.Username).Plugins
	}
	err := profile.FeedUserProfile(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	log.InfoContextf(request.Ctx, "UserProfile")
	err = recall.FeedRecall(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	log.InfoContextf(request.Ctx, "Recall")
	err = rank.FeedRank(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	log.InfoContextf(request.Ctx, "Rank")
	err = rerank.FeedRerank(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	return nil
}
