package rerank

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/feed/rerank/utils"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
)

const (
	minRetNum = 30
	maxRetNum = 50
)

func FeedRerank(request *model.FeedRequest) error {
	log.InfoContext(request.Ctx, "Rerank")
	existSet := set.NewStringSet(100)
	err := fillUserRead(request)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
	}
	chain := utils.NewExistFilterChain(existSet, request.UserRead)
	ret := make([]model.Item, 0, maxRetNum)
	rand.Seed(time.Now().UnixNano())
	feedNum := rand.Intn(maxRetNum-minRetNum) + minRetNum
	for k := 0; k < feedNum; k++ {
		var fristItem *model.Item
		for i := range request.RankResults {
			item := request.RankResults[i]
			if !chain.DoFilter(&item) {
				continue
			}
			if fristItem == nil {
				fristItem = &item
			}
			if !creatorCanInsert(ret[common.MaxInt(0, k-5):k], &item) {
				continue
			}
			fristItem = &item
			break
		}
		if fristItem != nil {
			ret = append(ret, *fristItem)
			existSet.Add(strconv.Itoa(fristItem.Music.ID))
		}
	}
	request.FeedResults = ret
	request.RankResults = nil
	return nil
}

func fillUserRead(request *model.FeedRequest) error {
	dbUserRead, err := service.PubUserReadService.GetUserRead(request.Ctx, request.Username, request.AnyUser)
	if err != nil {
		log.ErrorContext(request.Ctx, err)
		return err
	}
	if dbUserRead == nil {
		return nil
	}
	musicRead := strings.Split(dbUserRead.MusicRead, model.MusicReadDelimter)
	userRead := set.NewStringSet(len(musicRead) * 2)
	userRead.Add(musicRead...)
	request.UserRead = userRead
	return nil
}

func creatorCanInsert(items []model.Item, item *model.Item) bool {
	if len(items) <= 0 {
		return true
	}
	if item == nil {
		return false
	}
	if item.Music.ID <= 0 {
		return false
	}
	if item.Music.CreatorIDs == "" {
		return false
	}
	creators := strings.Split(item.Music.CreatorIDs, model.CreatorDelimiter)
	for i := range creators {
		creator := creators[i]
		for j := range items {
			if strings.Index(items[j].Music.CreatorIDs, creator) != -1 {
				return false
			}
		}
	}
	return true
}
