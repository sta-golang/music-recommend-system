package feed

import (
	"fmt"
	"time"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/feed/profile"
	"github.com/sta-golang/music-recommend/feed/rank"
	"github.com/sta-golang/music-recommend/feed/recall"
	"github.com/sta-golang/music-recommend/feed/rerank"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
)

var (
	requestNilErr = fmt.Errorf("request not be nil")
)

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
	log.InfoContextf(request.Ctx, "UserProfile end")
	err = recall.FeedRecall(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	log.InfoContextf(request.Ctx, "Recall end")
	err = rank.FeedRank(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	log.InfoContextf(request.Ctx, "Rank end")
	err = rerank.FeedRerank(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	log.InfoContextf(request.Ctx, "Rerank end")
	go func() {
		err := service.PubUserReadService.AddUserRead(request.Ctx, request.Username, request.FeedResults, request.AnyUser)
		if err != nil {
			log.ErrorContext(request.Ctx, err)
		}
	}()
	time.Sleep(time.Millisecond * 30)
	return nil
}
