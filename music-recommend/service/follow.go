package service

import (
	"errors"
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
	"strconv"
)

const (
	followCacheKey = "follow_%s_f"
	unFollowCacheKey = "unFollow_%s_f"

	followExpireTime = cache.Hour * 24

	followLimit = 20
)

type followService struct {
}


var PubFollowService = followService{}

// Follow 关注一个作者
// 这里使用内存的存储结构的不太好 一开始是直接用数组 增加O(1) 删除O(n)
// 理想情况下 ： 这里的数据结构应该是增删都是O(1)
// 然后还支持顺序读 所以这里比较麻烦
// 如果是支持这样的话
// 如果增删复杂度是O(lgn) 可以用红黑树

// 目前是使用一个数组加一个hash表 数组存关注的 hash表存取关的
// 但是还是有问题 可能 内存放不下。因为这个数据很多的 目前4000个作者。
// 极端情况 作者 4000 全部关注的情况下就是 100000个极端用户 就会用掉3g内存 还是很大的开销的
// 所以实际的实现起来还是 别全部都放入内存。这里可以优化 这里主要为了方便
// 如果正常的情况 都不太好做 因为你新关注就要放到队头。取的时候好取
// 但是 取10-20就不好取了。
// 一种是保存队头的方法 第0号一个队头。 第 10号一个队头 然后还得记录前队头新增了多少
// 如果队头插了一个这个时候队头新增数量就是1 要取10-20的时候 拿到第10位置的队头 看下前面插入了多少所以知道带跳页就很麻烦
// 所以最简单方式就是不让跳页 从头遍历到尾。取多少个就是多少个。 这块有好的idea可以和我探讨一下
// 这块我就偷懒去做了
func (fs *followService) Follow(creatorID int, username string) *er.Error {
	if _, ok := cache.PubCacheService.Get(fmt.Sprintf(creatorDetailCacheFmt, creatorID));!ok {
		return er.NewError(common.NotFound, common.CreatorNotExistErr)
	}

	affected, err := model.NewFollowCreatorMysql().Insert(&model.FollowCreator{
		CreatorID: creatorID,
		Username:  username,
	})
	if err != nil {
		log.Error(err)
		return er.NewError(common.DBCreateErr, err)
	}
	if !affected {
		return er.NewError(common.Success, common.FollowRepeatWarn)
	}
	key := fmt.Sprintf(followCacheKey, username)
	if val, ok := cache.PubCacheService.Get(key);ok {
		arr := val.([]int)
		arr = append(arr, creatorID)
		return nil
	}
	arr := make([]int, 0)
	arr = append(arr, creatorID)
	fs.setFollowCache(username, arr)
	return nil
}

func (fs *followService) setFollowCache(username string, creators []int) {
	priority := cache.Six
	key := fmt.Sprintf(followCacheKey, username)
	user, _ := PubUserService.meInfo(username)
	priority += cache.Priority(user.LastMonthLoginNum / 10) * cache.One
	cache.PubCacheService.Set(key, creators, cache.Hour * 24, priority)
}

func (fs *followService) UnFollow(creatorID int, username string) *er.Error {
	if _, ok := cache.PubCacheService.Get(fmt.Sprintf(creatorDetailCacheFmt, creatorID));!ok {
		return er.NewError(common.NotFound, common.CreatorNotExistErr)
	}
	affected, err := model.NewFollowCreatorMysql().Delete(username, creatorID)
	if err != nil {
		log.Error(err)
		return er.NewError(common.DBCreateErr, err)
	}
	if !affected {
		return er.NewError(common.Success, common.UnFollowRepeatWarn)
	}
	key := fmt.Sprintf(unFollowCacheKey, username)
	if val, ok := cache.PubCacheService.Get(key); ok {
		unFollowSet := val.(*set.HashSet)
		unFollowSet.Add(creatorID)
		return nil
	}
	unFollowSet := set.NewHashSet()
	unFollowSet.Add(strconv.Itoa(creatorID))
	// 它一定比follow的晚过期
	cache.PubCacheService.Set(key, unFollowSet, followExpireTime, cache.Ten)
	return nil
}

func (fs *followService) FollowList(username string, page int) ([]*model.Creator, *er.Error) {
	key := fmt.Sprintf(followCacheKey, username)
	if val, ok := cache.PubCacheService.Get(key); ok {
		creatorIDs := val.([]int)
		return fs.doGetFollowList(username, creatorIDs, page), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ids, err := model.NewFollowCreatorMysql().SelectFollows(username, 0, 9999)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		fs.setFollowCache(username, ids)
		return ids, err
	})
	if err != nil {
		return nil, er.NewError(common.DBFindErr, err)
	}
	if ret != nil {
		return fs.doGetFollowList(username, ret.([]int), page), nil
	}
	return nil, nil
}

func (fs *followService) doGetFollowList(username string, creatorIDs []int, page int) []*model.Creator {
	var unFollowSet *set.HashSet
	if val, ok := cache.PubCacheService.Get(fmt.Sprintf(unFollowCacheKey, username));ok {
		unFollowSet = val.(*set.HashSet)
	}
	start := common.MaxInt(0, len(creatorIDs) - 1 - followLimit * (page - 1))
	cnt := 0
	ret := make([]*model.Creator, 0, common.MinInt(20, start))
	for i := start; i >= 0 && cnt < followLimit; i-- {
		if unFollowSet != nil && unFollowSet.Contains(creatorIDs[i]) {
			continue
		}
		creator, err := PubCreatorService.getCreatorWithCache(creatorIDs[i])
		if creator == nil || err != nil {
			continue
		}
		ret = append(ret, creator)
		cnt++
	}
}