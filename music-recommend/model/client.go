package model

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/music-recommend/db"
)

const (
	dbMusicRecommendNameTest = "music_recommend_test"
)

var noResultErr = sql.ErrNoRows

func client(dbName string) *sqlx.DB {
	return db.GetDB(dbName)
}
