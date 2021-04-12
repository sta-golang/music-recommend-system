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
	playlistUserKey = "playlist_%d_all"
)

type playlistService struct {
}

var PubPlaylistService = &playlistService{}

func (ps *playlistService) GetPlaylistForUser(ctx context.Context, userID int) ([]model.Playlist, error) {
	playlists, err := model.NewPlaylistMysql().SelectForUser(ctx, userID)
	if err != nil {
		log.Error(err)
		return nil, common.ServerErr
	}
	return playlists, nil
}

func (ps *playlistService) GetPlaylistForUserWithCache(ctx context.Context, userID int) ([]model.Playlist, error) {
	key := fmt.Sprintf(playlistUserKey, userID)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Playlist), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		retData, err := ps.GetPlaylistForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		if retData == nil {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, retData, cache.Hour*25, cache.Seven)
		return retData, nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.([]model.Playlist), nil
}

func (ps *playlistService) AddPlaylistForUserWithCache(ctx context.Context, name string, userID int) error {
	playlist := model.Playlist{
		UserID: userID,
		Name:   name,
	}
	affected, err := model.NewPlaylistMysql().Insert(ctx, &playlist)
	if err != nil {
		log.Error(err)
		return common.ServerErr
	}
	if !affected {
		return fmt.Errorf("歌单名:%s 已经存在", name)
	}
	key := fmt.Sprintf(playlistUserKey, userID)
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
	cache.PubCacheService.Set(key, cachePlaylist, cache.Hour*24, cache.Seven)
	return nil
}

func (ps *playlistService) GetMusicWithPlaylist(ctx context.Context) {

}

func (ps *playlistService) DeletePlaylistForUserWithCache(ctx context.Context, id, userID int) error {
	err := model.NewPlaylistMysql().DeleteMusicForPlaylist(ctx, id, userID)
	if err != nil {
		log.Error(err)
		return common.ServerErr
	}
	key := fmt.Sprintf(playlistUserKey, userID)
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
	}
	return nil
}
