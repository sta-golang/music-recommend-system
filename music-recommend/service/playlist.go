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
	playlistMusicKey  = "playlist_%d_music"
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

func (ps *playlistService) GetPlaylistMusicWithCache(ctx context.Context, playlistID int) ([]model.Music, error) {
	key := fmt.Sprintf(playlistMusicKey, playlistID)
	if val, ok := cache.PubCacheService.Get(key); ok {
		if val == nil {
			return nil, nil
		}
		return val.([]model.Music), nil
	}
	retData, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		ret, err := ps.GetPlaylistMusic(ctx, playlistID)
		if err != nil {
			return nil, err
		}
		if len(ret) <= 0 {
			cache.PubCacheService.Set(key, nil, cache.Hour, cache.One)
			return nil, nil
		}
		cache.PubCacheService.Set(key, ret, cache.Hour*4, cache.Four)
		return ret, nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	if retData == nil {
		return nil, nil
	}
	return retData.([]model.Music), nil
}

func (ps *playlistService) userHasModifyAuthor(ctx context.Context, plyalistID int, username string) (bool, error) {
	playlist, err := ps.GetPlaylistDetailWithCache(ctx, plyalistID)
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}
	if playlist == nil {
		return false, fmt.Errorf("该歌单不存在")
	}
	return playlist.Username == username, nil
}

func (ps *playlistService) AddMusicToPlaylist(ctx context.Context, musicID, playlistID int, username string) error {
	music, sErr := PubMusicService.GetMusic(musicID)
	if sErr != nil && sErr.Err != nil {
		if sErr.Code != common.NotFound {
			log.ErrorContext(ctx, sErr)
			return sErr.Err
		}
	}
	if music == nil {
		return fmt.Errorf("该音乐不存在")
	}
	if ok, err := ps.userHasModifyAuthor(ctx, playlistID, username); !ok || err != nil {
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		if !ok {
			return fmt.Errorf("你没有修改此歌单的权限")
		}
	}
	cachePlaylistMusic, err := ps.GetPlaylistMusicWithCache(ctx, playlistID)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	if len(cachePlaylistMusic) > model.MaxPlaylistMusicSize {
		return fmt.Errorf("同一个歌单最多添加%d首歌曲", model.MaxPlaylistMusicSize)
	}
	affected, err := model.NewPlaylistMysql().AddMusicForPlaylist(ctx, musicID, playlistID, username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	if !affected {
		return fmt.Errorf("添加失败 无权限或者歌曲已在歌单中")
	}
	key := fmt.Sprintf(playlistMusicKey, playlistID)
	if cachePlaylistMusic != nil {
		oldPlaylistMusic := cachePlaylistMusic
		endPlaylistMusic := oldPlaylistMusic[len(oldPlaylistMusic)-1]
		for i := len(oldPlaylistMusic) - 1; i > 0; i-- {
			oldPlaylistMusic[i] = oldPlaylistMusic[i-1]
		}
		oldPlaylistMusic[0] = *music
		oldPlaylistMusic = append(oldPlaylistMusic, endPlaylistMusic)
		cachePlaylistMusic = oldPlaylistMusic
	}
	if len(cachePlaylistMusic) <= 0 {
		cachePlaylistMusic = []model.Music{*music}
	}
	cache.PubCacheService.Set(key, cachePlaylistMusic, cache.Hour*2, cache.Four)
	return nil
}

func (ps *playlistService) DeleteMusicForPlaylist(ctx context.Context, musicID, playlistID int, username string) error {
	if ok, err := ps.userHasModifyAuthor(ctx, playlistID, username); !ok || err != nil {
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		if !ok {
			return fmt.Errorf("你没有修改此歌单的权限")
		}
	}
	affected, err := model.NewPlaylistMysql().DeleteMusicForPlaylist(ctx, musicID, playlistID, username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	if !affected {
		return fmt.Errorf("无权限或者该音乐不在此表单中")
	}
	index := -1
	key := fmt.Sprintf(playlistMusicKey, playlistID)
	if val, ok := cache.PubCacheService.Get(key); ok && val != nil {
		playlistMusics := val.([]model.Music)
		for i := range playlistMusics {
			if playlistMusics[i].ID == musicID {
				index = i
				break
			}
		}
		if index == -1 {
			return nil
		}
		for i := index; i < len(playlistMusics)-1; i++ {
			playlistMusics[i] = playlistMusics[i+1]
		}
		playlistMusics = playlistMusics[:len(playlistMusics)-1]
		cache.PubCacheService.Set(key, playlistMusics, cache.Hour*2, cache.Four)
	}
	return nil
}

// 这里前端没有传值 就一次获取把。 后续想改的时候可以加page页
func (ps *playlistService) GetPlaylistMusic(ctx context.Context, playlistID int) ([]model.Music, error) {
	musics, err := model.NewPlaylistMysql().SelectMusicsForPlaylist(ctx, playlistID, 0, 5000)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return musics, nil
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
