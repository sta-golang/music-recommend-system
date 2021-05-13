package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
)

const (
	searchKeyworld  = "searchKeyworld_%s"
	searchMusic     = "searchMusic_%s"
	maxSearchConunt = 100
)

type searchService struct {
}

var PubSearchService = &searchService{}

// 首先搜索名字完全一样的一个作者 然后搜索名字像这个的作者的前10个
// 然后one-hot编码计算相似度取前3个作者
// 然后搜索名字完全一样的两首歌曲 然后搜索像的三首歌曲
// 然后通过前三个作者拿到它们的热门歌曲的 1：3 2：1 3：1
func (ss *searchService) SearchKeyworld(ctx context.Context, keyworld string) (*model.SearchResult, error) {
	ret, err := common.SingleRunGroup.Do(fmt.Sprintf(searchKeyworld, keyworld), func() (interface{}, error) {
		graph := ss.buildDag(ctx, keyworld)
		defer graph.DestoryAsync()
		graph.Do(ctx, false)
		ret, err := graph.GetRootTask().GetRet()
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, nil
		}
		return ret, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*model.SearchResult), nil
}

func (ss *searchService) SearchMusics(ctx context.Context, keyworld string) ([]model.Music, error) {
	ret, err := common.SingleRunGroup.Do(fmt.Sprintf(searchMusic, keyworld), func() (interface{}, error) {
		var ret []model.Music
		musics, err := model.NewSearchMysql().SearchForMusics(ctx, keyworld, 0, 3)
		if err != nil {
			return nil, err
		}
		var mSet *set.StringSet
		if len(musics) > 0 {
			mSet = set.NewStringSet(len(musics) * 2)
			for i := range musics {
				mSet.Add(strconv.Itoa(musics[i].ID))
				ret = append(ret, musics[i])
			}
		}
		musics, err = model.NewSearchMysql().SearchForMusicsLike(ctx, keyworld, 0, maxSearchConunt)
		if err != nil {
			return nil, err
		}
		for i := range musics {
			if mSet.Contains(strconv.Itoa(musics[i].ID)) {
				continue
			}
			ret = append(ret, musics[i])
		}
		return ret, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.([]model.Music), nil
}

func (ss *searchService) buildDag(ctx context.Context, keyworld string) *dag.DagTasks {
	keyOneHot := []rune(keyworld)
	creatorMerge := dag.NewTask("creatorMerge", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		var ret []model.Creator
		mergeSet := set.NewStringSet(10)
		for k := 0; k < helper.GetSubTaskSize(); k++ {
			subRet, err := helper.GetSubTaskRetForIndex(k)
			if err != nil {
				continue
			}
			if subRet == nil {
				continue
			}
			subMusics := subRet.([]model.Creator)
			if len(subMusics) <= 0 {
				continue
			}
			for c := range subMusics {
				if mergeSet.Contains(strconv.Itoa(subMusics[c].ID)) {
					continue
				}
				mergeSet.Add(strconv.Itoa(subMusics[c].ID))
				ret = append(ret, subMusics[c])
			}
		}
		sort.Slice(ret, func(i, j int) bool {
			if ss.oneHot(keyOneHot, []rune(ret[i].Name)) > ss.oneHot(keyOneHot, []rune(ret[j].Name)) {
				return true
			}
			return false
		})
		ret = ret[:common.MinInt(3, len(ret))]
		return ret, nil
	})
	creatorMerge.AddSubTask(dag.NewTask("SearchForCreator", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		ret, err := model.NewSearchMysql().SearchForCreator(ctx, keyworld, 0, 1)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
		return ret, nil
	}))
	creatorMerge.AddSubTask(dag.NewTask("SearchForCreatorLike", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		ret, err := model.NewSearchMysql().SearchForCreatorLike(ctx, keyworld, 0, 10)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
		return ret, nil
	}))
	// --------------------------------- music
	musicSearchMerge := dag.NewTask("musicSearchMerge", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		var ret []model.Music
		subRet, err := helper.GetSubTaskRet("musicSearch")
		if err == nil && subRet != nil {
			subMusics := subRet.([]model.Music)
			ret = append(ret, subMusics...)
		}
		subRet, err = helper.GetSubTaskRet("musicSearchLike")
		if err == nil && subRet != nil {
			subMusics := subRet.([]model.Music)
			ret = append(ret, subMusics...)
		}
		return ret, nil
	})
	musicSearchMerge.AddSubTask(dag.NewTask("musicSearch", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		ret, err := model.NewSearchMysql().SearchForMusics(ctx, keyworld, 0, 2)
		if err != nil {
			log.ErrorContext(ctx, err)
			return nil, err
		}
		return ret, nil
	}))
	musicSearchMerge.AddSubTask(dag.NewTask("musicSearchLike", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		ret, err := model.NewSearchMysql().SearchForMusicsLike(ctx, keyworld, 0, 3)
		if err != nil {
			log.ErrorContext(ctx, err)
			return nil, err
		}
		return ret, nil
	}))
	// ------------------------------------------- merge
	merge := dag.NewTask("merge", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		ret := &model.SearchResult{}
		subRet, err := helper.GetSubTaskRet("creatorMerge")
		if err == nil && subRet != nil {
			subCreators := subRet.([]model.Creator)
			ret.Creators = subCreators
		}
		subRet, err = helper.GetSubTaskRet("musicSearchMerge")
		if err == nil && subRet != nil {
			subMusics := subRet.([]model.Music)
			ret.Musics = subMusics
		}
		if len(ret.Creators) >= 1 {
			ms, _ := PubMusicService.GetMusicForCreator(ret.Creators[0].ID, 1)
			if len(ms) >= 0 {
				ret.Musics = append(ret.Musics, ms[:common.MinInt(3, len(ms))]...)
			}
		}
		if len(ret.Creators) >= 2 {
			ms, _ := PubMusicService.GetMusicForCreator(ret.Creators[1].ID, 1)
			if len(ms) >= 0 {
				ret.Musics = append(ret.Musics, ms[:common.MinInt(1, len(ms))]...)
			}
		}
		if len(ret.Creators) >= 3 {
			ms, _ := PubMusicService.GetMusicForCreator(ret.Creators[2].ID, 1)
			if len(ms) >= 0 {
				ret.Musics = append(ret.Musics, ms[:common.MinInt(1, len(ms))]...)
			}
		}
		// 去重
		musicSet := set.NewHashSet(10)
		musics := make([]model.Music, 0, len(ret.Musics))
		for i := range ret.Musics {
			if musicSet.Contains(strconv.Itoa(ret.Musics[i].ID)) {
				continue
			}
			musicSet.Add(strconv.Itoa(ret.Musics[i].ID))
			musics = append(musics, ret.Musics[i])
		}
		ret.Musics = musics
		return ret, nil
	})
	merge.AddSubTask(musicSearchMerge)
	merge.AddSubTask(creatorMerge)
	return dag.NewDag(merge)
}

func (ss *searchService) oneHot(one, two []rune) float64 {
	keySet := set.NewHashSet(len(one))
	cnt := 0
	sum := 0
	for _, r := range one {
		if keySet.Contains(r) {
			continue
		}
		keySet.Add(r)
		sum++
	}
	for _, r := range two {
		if keySet.Contains(r) {
			cnt++
		} else {
			sum++
		}
	}
	return float64(cnt) / float64(sum)
}
