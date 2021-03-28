package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/source"
	"github.com/sta-golang/music-recommend/config"
)

type sqlxDB struct {
	dbMap map[string]*sqlx.DB
}

const (
	dbFmt      = "%s:%s@tcp(%s)/%s?%s"
	sourceName = "database"
)

func (sd *sqlxDB) Name() string {
	return sourceName
}

func (sd *sqlxDB) Sync() error {
	for _, val := range sd.dbMap {
		err := val.Close()
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

var globalDB *sqlxDB

// InitDB 初始化数据库 连接数据库使用
func InitDB() error {
	globalDB = &sqlxDB{dbMap: make(map[string]*sqlx.DB, 1)}
	source.Monitoring(globalDB)
	var err error
	cfg := config.GlobalConfig()
	for key, val := range cfg.DBConfigs {
		ipPort := val.Target
		if ipPort == "" {
			ipPort = "127.0.0.1:3306"
		}
		globalDB.dbMap[key], err = sqlx.Open(val.DriverName, fmt.Sprintf(dbFmt,
			val.UserName, val.PassWord, val.Target, val.DBName, val.Args))
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func GetDB(key string) *sqlx.DB {
	return globalDB.dbMap[key]
}
