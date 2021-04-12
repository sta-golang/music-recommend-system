package model

import (
	"context"
	"fmt"

	"github.com/sta-golang/go-lib-utils/log"

	tm "github.com/sta-golang/go-lib-utils/time"
)

type Playlist struct {
	ID          int     `json:"id" db:"id"`
	Status      int32   `json:"status" db:"status"`
	UserID      int     `json:"user_id" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	ImageUrl    string  `json:"image_url" db:"image_url"`
	HotScore    float64 `json:"hot_score" db:"hot_score"`
	TagNames    string  `json:"tag_names" db:"tag_names"`
	Description string  `json:"description" db:"description"`
	CreateTime  string  `json:"create_time" db:"create_time"`
	UpdateTime  string  `json:"update_time" db:"update_time"`
}

const (
	tablePlaylist      = "playlist"
	tablePlaylistMusic = "playlist_music"
)

type playlistMysql struct {
}

var oncePlaylistMysql = playlistMysql{}

func NewPlaylistMysql() *playlistMysql {
	return &oncePlaylistMysql
}

func (pm *playlistMysql) Insert(ctx context.Context, p *Playlist) (bool, error) {
	sql := fmt.Sprintf("insert ignore into %s(status,user_id,name,image_url,create_time,update_time) values(?,?,?,?,?,?)", tablePlaylist)
	res, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, 0, p.UserID, p.Name, p.ImageUrl, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)

		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return false, err
	}
	insertID, err := res.LastInsertId()
	if err != nil {
		log.Error(err)
		return true, err
	}
	p.ID = int(insertID)
	return rows > 0, err
}

func (pm *playlistMysql) Update(ctx context.Context, p *Playlist) error {
	sql := fmt.Sprintf("update %s set status=?,image_url=?,hot_score=?,tag_names=?,description=?,update_time=? where id = ?", tablePlaylist)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, p.Status, p.ImageUrl, p.HotScore, p.TagNames, p.Description, tm.GetNowDateTimeStr(), p.ID)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (pm *playlistMysql) SelectForUser(ctx context.Context, userID int) (playlists []Playlist, err error) {
	sql := fmt.Sprintf("select id, name,user_id from %s where user_id = ? order by id desc", tablePlaylist)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &playlists, sql, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return
}

func (pm *playlistMysql) SelectMusicsForPlaylist(ctx context.Context, id, userID, pos, limit int) (musics []Music, err error) {
	sql := fmt.Sprintf("select * from %s where id in (select music_id from %s where playlist_id = ? and user_id = ? order by id desc limit ?, ?)", tableMusic, tablePlaylistMusic)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &musics, sql, id, userID, pos, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}

// AddMusicForPlaylist 和Delete说明类似
func (pm *playlistMysql) AddMusicForPlaylist(ctx context.Context, musicID, playlistID, userID int) (bool, error) {
	sql := fmt.Sprintf("insert ignore into %s(music_id, playlist_id,user_id, create_time, update_time) values(?,?,?,?,?)", tablePlaylistMusic)
	res, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, musicID, playlistID, userID, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr())
	if err != nil {
		log.Error(err)

		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return false, err
	}
	return rows > 0, nil
}

// DeleteMusicForPlaylist 传入userID是因为可以防止别人通过请求直接删除而不校验权限
func (pm *playlistMysql) DeleteMusicForPlaylist(ctx context.Context, id, userID int) error {
	sql := fmt.Sprintf("delete from %s where id = ? and user_id = ?", tablePlaylistMusic)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, id, userID)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (pm *playlistMysql) Select(ctx context.Context, id int) (*Playlist, error) {
	var ret Playlist
	sql := fmt.Sprintf("select * from %s where id = ?", tablePlaylist)
	err := client(dbMusicRecommendNameTest).GetContext(ctx, &ret, sql, id)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &ret, nil
}
