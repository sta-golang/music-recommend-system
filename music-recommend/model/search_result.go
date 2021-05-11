// Package model provides ...
package model

import (
	"context"
	"fmt"

	"github.com/sta-golang/go-lib-utils/log"
)

type SearchResult struct {
	Musics   []Music   `json:"songs"`
	Creators []Creator `json:"artists"`
}

type searchMysql struct {
}

var onceSearchMysql = searchMysql{}

func NewSearchMysql() *searchMysql {
	return &onceSearchMysql
}

func (sm *searchMysql) SearchForCreator(ctx context.Context, keyword string, pos, limit int) (creators []Creator, err error) {
	sql := fmt.Sprintf("select * from %s where name = ? limit ?, ?", tableCreator)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &creators, sql, keyword, pos, limit)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return
}

func (sm *searchMysql) SearchForCreatorLike(ctx context.Context, keyword string, pos, limit int) (creatos []Creator, err error) {
	sql := fmt.Sprintf("select * from %s where name like '%%%s%%' order by hot_score desc limit ?,?", tableCreator, keyword)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &creatos, sql, pos, limit)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return
}

func (sm *searchMysql) SearchForMusics(ctx context.Context, keyword string, pos, limit int) (musics []Music, err error) {
	sql := fmt.Sprintf("select * from %s where name = ? and status = ? limit ?,?", tableMusic)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &musics, sql, keyword, StatusLoadMusicFinish, pos, limit)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return
}

func (sm *searchMysql) SearchForMusicsLike(ctx context.Context, keyword string, pos, limit int) (musics []Music, err error) {
	sql := fmt.Sprintf("select * from %s where status = ? and name like '%%%s%%' limit ?, ?", tableMusic, keyword)
	err = client(dbMusicRecommendNameTest).SelectContext(ctx, &musics, sql, StatusLoadMusicFinish, pos, limit)
	if err != nil {
		log.ErrorContext(ctx, err)
		return nil, err
	}
	return
}
