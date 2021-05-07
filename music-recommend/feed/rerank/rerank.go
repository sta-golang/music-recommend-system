package rerank

import (
	"math/rand"
	"strings"
	"time"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/feed/rerank/utils"
	"github.com/sta-golang/music-recommend/model"
)

const (
	minRetNum = 30
	maxRetNum = 50
)

func FeedRerank(request *model.FeedRequest) error {
	existSet := set.NewHashSet(maxRetNum << 1)
	chain := utils.NewExistFilterChain(existSet)
	var fristItem *model.Item
	ret := make([]model.Item, 0, maxRetNum)
	rand.Seed(time.Now().UnixNano())
	feedNum := rand.Intn(maxRetNum-minRetNum) + minRetNum
	for k := 0; k < feedNum; k++ {
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
		ret = append(ret, *fristItem)
		existSet.Add(fristItem.Music.ID)
	}
	request.FeedResults = ret
	request.RankResults = nil
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
