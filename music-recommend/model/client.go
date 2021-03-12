package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/music-recommend/db"
)

const (
	dbMusicRecommendNameTest = "music_recommend_test"
)

func client(dbName string) *sqlx.DB {
	return db.GetDB(dbName)
}
