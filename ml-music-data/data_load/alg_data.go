package data_load

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/model"
)

type algData struct {
	writer io.Writer
	table  map[int]bool
}

func NewAlgData(writer io.Writer) *algData {
	return &algData{
		writer: writer,
	}
}

func (ad *algData) WriterMusic2Vec() error {
	return ad.getPlaylistVec()
}

func (ad *algData) getPlaylistVec() error {
	cl := WangYiYunCrawler{}
	playlists, err := cl.GetPlayListIDs()
	if err != nil {
		log.ConsoleLogger.Error(err)
		return err
	}
	log.ConsoleLogger.Info("crawler playlist finish!")
	ad.table = make(map[int]bool, 4096)
	for cnt, playlist := range playlists {
		log.ConsoleLogger.Infof("load playlist once! score is %.2f", float64(cnt*100)/float64(len(playlists)))
		res, err := cl.CrawlerPlaylistsDetail(playlist)
		if err != nil {
			log.ConsoleLogger.Error(err)
			return err
		}
		strArr := make([]string, 0, len(res.Playlist.SongIDs))
		for _, song := range res.Playlist.SongIDs {
			if !ad.canAdd(song.ID) {
				continue
			}
			strArr = append(strArr, strconv.Itoa(song.ID))
		}
		if len(strArr) <= 0 {
			log.ConsoleLogger.Warn("length <= 0")
			continue
		}
		rand.Shuffle(len(strArr), func(i, j int) {
			strArr[i], strArr[j] = strArr[j], strArr[i]
		})
		ss := fmt.Sprintf("%s\n", strings.Join(strArr, " "))
		ad.writer.Write(str.StringToBytes(&ss))
		if err != nil {
			log.ConsoleLogger.Error(err)
			continue
		}
	}

	return nil
}

func (ad *algData) canAdd(songID int) bool {
	if val, ok := ad.table[songID]; ok {
		return val
	}
	music, err := model.NewMusicMysql().SelectByID(songID)
	if err != nil {
		log.ConsoleLogger.Error(songID, " SelectMysql err : ", err)
		return false
	}
	if music != nil && music.Status == model.MusicHasMusicUrlStatus {
		ad.table[songID] = true
		return true
	}
	ad.table[songID] = false
	return false
}
