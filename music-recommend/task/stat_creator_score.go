package task

import (
	"context"
	"time"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/model"
)

type CreatorHotScoreTask struct {
}

func NewCreatorHotScoreTask() TaskHandler {
	return &CreatorHotScoreTask{}
}

func (ch *CreatorHotScoreTask) TaskTimeInterval() time.Duration {
	return time.Hour * 24
}

func (ch *CreatorHotScoreTask) RunNow() bool {
	return true
}

func (ch *CreatorHotScoreTask) Handler(ctx context.Context) error {
	scors, err := model.NewCreatorMysql().GetScoresTable(ctx)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	limit := 200
	start, end := 0, limit

	for {
		if end > len(scors) {
			end = len(scors)
		}
		if start == len(scors) {
			break
		}
		err = model.NewCreatorMysql().UpdateCreatorScores(ctx, scors[start:end])
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		start = end
		end += limit
	}
	return nil
}

func (ch *CreatorHotScoreTask) Name() string {
	return "CreatorHotScoreTask"
}
