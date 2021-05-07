package model

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"

	tm "github.com/sta-golang/go-lib-utils/time"
)

type Playlist struct {
	ID          int     `json:"id" db:"id"`
	Status      int32   `json:"status" db:"status"`
	Username    string  `json:"username" db:"username"`
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

	MaxPlaylistMusicSize = 200
)

type playlistMysql struct {
}

var oncePlaylistMysql = playlistMysql{}

func NewPlaylistMysql() *playlistMysql {
	return &oncePlaylistMysql
}

func (pm *playlistMysql) Insert(ctx context.Context, p *Playlist) (bool, error) {
	sql := fmt.Sprintf("insert ignore into %s(status,username,name,image_url,description,create_time,update_time) values(?,?,?,?,?,?,?)", tablePlaylist)
	res, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, 0, p.Username, p.Name, p.ImageUrl, p.Description, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr())
	if err != nil {
		log.ErrorContext(ctx, err)

		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}
	insertID, err := res.LastInsertId()
	if err != nil {
		log.ErrorContext(ctx, err)
		return true, err
	}
	p.ID = int(insertID)
	return rows > 0, err
}

func (pm *playlistMysql) Update(ctx context.Context, p *Playlist) error {
	sql := fmt.Sprintf("update %s set status=?,image_url=?,hot_score=?,tag_names=?,description=?,update_time=? where id = ?", tablePlaylist)
	_, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, p.Status, p.ImageUrl, p.HotScore, p.TagNames, p.Description, tm.GetNowDateTimeStr(), p.ID)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	return err
}

func (pm *playlistMysql) SelectForUser(ctx context.Context, username string) (playlists []Playlist, err error) {
	sql := fmt.Sprintf("select id, name,username from %s where username = ? order by id desc", tablePlaylist)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &playlists, sql, username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return
}

func (pm *playlistMysql) SelectForHotScore(ctx context.Context, pos, limit int) (playlists []Playlist, err error) {
	sql := fmt.Sprintf("select id, name,username,image_url,hot_score from %s order by hot_score desc limit ?, ?", tablePlaylist)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &playlists, sql, pos, limit)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return
}

func (pm *playlistMysql) SelectMusicsForPlaylist(ctx context.Context, id, pos, limit int) (musics []Music, err error) {
	sql := fmt.Sprintf("select * from %s where id in (select music_id from %s where playlist_id = ? order by id desc ) limit ?, ?", tableMusic, tablePlaylistMusic)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &musics, sql, id, pos, limit)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return musics, nil
}

// AddMusicForPlaylist 和Delete说明类似
func (pm *playlistMysql) AddMusicForPlaylist(ctx context.Context, musicID, playlistID int, username string) (bool, error) {
	sql := fmt.Sprintf("insert ignore into %s(music_id, playlist_id,username, create_time, update_time) values(?,?,?,?,?)", tablePlaylistMusic)
	res, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, musicID, playlistID, username, tm.GetNowDateTimeStr(), tm.GetNowDateTimeStr())
	if err != nil {
		log.ErrorContext(ctx, err)

		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}
	return rows > 0, nil
}

func (pm *playlistMysql) DeletePlaylist(ctx context.Context, id int, username string) (bool, error) {
	affected := false
	err := newMysqlTransaction().Transaction(func(tx *sqlx.Tx) error {
		// TODO: 删除playlis单独封装一个函数  <13-04-21, FOUR SEASONS> //
		sql := fmt.Sprintf("delete from %s where id = ? and username = ?", tablePlaylist)
		res, err := tx.ExecContext(ctx, sql, id, username)
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			log.ErrorContext(ctx, err)
			return err
		}
		affected = rows > 0
		if !affected {
			return nil
		}
		// TODO: 删除playlist下所有音乐单独封装一个函数  <13-04-21, FOUR SEASONS> //
		sql = fmt.Sprintf("delete from %s where playlist_id = ? and username = ?", tablePlaylistMusic)
		_, err = tx.ExecContext(ctx, sql, id, username)
		if err != nil {
			log.ErrorContext(ctx, err)
		}
		return nil
	})
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}
	return affected, err

}

// DeleteMusicForPlaylist 传入userID是因为可以防止别人通过请求直接删除而不校验权限
func (pm *playlistMysql) DeleteMusicForPlaylist(ctx context.Context, musicID, playlistID int, username string) (bool, error) {
	sql := fmt.Sprintf("delete from %s where music_id = ? and playlist_id = ? and username = ?", tablePlaylistMusic)
	res, err := client(dbMusicRecommendNameTest).ExecContext(ctx, sql, musicID, playlistID, username)
	if err != nil {
		log.ErrorContext(ctx, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.ErrorContext(ctx, err)
		return false, err
	}

	return rows > 0, err
}

func (pm *playlistMysql) Select(ctx context.Context, id int) (*Playlist, error) {
	var ret Playlist
	sql := fmt.Sprintf("select * from %s where id = ?", tablePlaylist)
	err := client(dbMusicRecommendNameTest).GetContext(ctx, &ret, sql, id)
	if err == noResultErr {
		return nil, nil
	}
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return &ret, nil
}
