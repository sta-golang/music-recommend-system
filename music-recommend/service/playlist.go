package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
)

const (
	playlistUserKey   = "playlist_%s_all"
	playlistDetailKey = "playlist_%d_de"
	defDescription    = "暂无介绍"
	defImageUrl       = "https://gimg3.baidu.com/image_search/src=http%3A%2F%2Fwenhui.whb.cn%2Fu%2Fcms%2Fwww%2F201804%2F02165434mp7i.jpg&refer=http%3A%2F%2Fwenhui.whb.cn&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1620922142&t=cb795fafc42b7438723435d40553cba3 "
)

type playlistService struct {
}

var PubPlaylistService = &playlistService{}

func (ps *playlistService) GetPlaylistForUser(ctx context.Context, username string) ([]model.Playlist, error) {
	playlists, err := model.NewPlaylistMysql().SelectForUser(ctx, username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, common.ServerErr
	}
	return playlists, nil
}

func (ps *playlistService) GetPlaylistForUserWithCache(ctx context.Context, username string) ([]model.Playlist, error) {
	key := fmt.Sprintf(playlistUserKey, username)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Playlist), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		retData, err := ps.GetPlaylistForUser(ctx, username)
		if err != nil {
			return nil, err
		}
		if retData == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, retData, cache.Hour*12, cache.Three)
		return retData, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.([]model.Playlist), nil
}

func (ps *playlistService) AddPlaylistForUserWithCache(ctx context.Context, name, username string) error {
	playlist := model.Playlist{
		Username:    username,
		Name:        name,
		Description: defDescription,
		ImageUrl:    defImageUrl,
	}
	affected, err := model.NewPlaylistMysql().Insert(ctx, &playlist)
	if err != nil {
		log.ErrorContext(ctx, err)
		return common.ServerErr
	}
	if !affected {
		return fmt.Errorf("歌单名:%s 已经存在", name)
	}
	key := fmt.Sprintf(playlistUserKey, username)
	var cachePlaylist []model.Playlist
	if val, ok := cache.PubCacheService.Get(key); ok && val != nil {
		oldPlaylist := val.([]model.Playlist)
		// 这种情况理论上不应该出现
		if len(oldPlaylist) > 0 && playlist.ID < oldPlaylist[0].ID {
			cache.PubCacheService.Delete(key)
			return nil
		}
		endPlaylist := oldPlaylist[len(oldPlaylist)-1]
		for i := len(oldPlaylist) - 1; i > 0; i-- {
			oldPlaylist[i] = oldPlaylist[i-1]
		}
		oldPlaylist[0] = playlist
		oldPlaylist = append(oldPlaylist, endPlaylist)
		cachePlaylist = oldPlaylist
	}
	if len(cachePlaylist) <= 0 {
		cachePlaylist = []model.Playlist{playlist}
	}
	cache.PubCacheService.Set(key, cachePlaylist, cache.Hour*8, cache.Three)
	return nil
}

func (ps *playlistService) GetPlaylistDetailWithCache(ctx context.Context, id int) (*model.Playlist, error) {
	key := fmt.Sprintf(playlistDetailKey, id)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, fmt.Errorf(common.NotFoundMessage)
		}
		return val.(*model.Playlist), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ret, err := ps.GetPlaylistDetail(ctx, id)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, fmt.Errorf(common.NotFoundMessage)
		}
		return ret, nil
	})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*model.Playlist), nil
}

func (ps *playlistService) GetPlaylistDetail(ctx context.Context, id int) (*model.Playlist, error) {
	ret, err := model.NewPlaylistMysql().Select(ctx, id)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (ps *playlistService) GetMusicWithPlaylist(ctx context.Context) {

}

func (ps *playlistService) DeletePlaylistForUserWithCache(ctx context.Context, id int, username string) error {
	affected, err := model.NewPlaylistMysql().DeletePlaylist(ctx, id, username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return common.ServerErr
	}
	if !affected {
		return nil
	}
	key := fmt.Sprintf(playlistUserKey, username)
	if val, ok := cache.PubCacheService.Get(key); ok && val != nil {
		playlists := val.([]model.Playlist)
		// 二分查找
		index := sort.Search(len(playlists), func(i int) bool {
			return playlists[i].ID <= id
		})
		if index >= len(playlists) || playlists[index].ID != id {
			return nil
		}
		for i := index; i < len(playlists)-1; i++ {
			playlists[i] = playlists[i+1]
		}
		playlists = playlists[:len(playlists)-1]
		cache.PubCacheService.Set(key, playlists, cache.Hour*8, cache.Three)
	}
	return nil
}
