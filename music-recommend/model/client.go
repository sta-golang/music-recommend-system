package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/music-recommend/db"
)

func client(dbName string) *sqlx.DB {
	return db.GetDB(dbName)
}
