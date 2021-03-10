package data_load

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/cmd"
	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/sta-golang/music-recommend/model"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	pythonCmd = "python"
	crawlerPythonPath = "./music_download/crawler.py"
	playlistBaseUrlFmt = "https://music.163.com/api/playlist/detail?id=%s"
)

type WangYiYunCrawler struct {
}

func (wc *WangYiYunCrawler) crawlerInfo() ([]byte, error) {
	command, err := cmd.ExecCmd(pythonCmd, crawlerPythonPath)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	message := command.OutMessage
	if command.RunErr != nil {
		log.Warn(command.RunErr)
	}
	if len(message) > 0 {
		return message, nil
	}
	return nil, nil
}

func (wc *WangYiYunCrawler) GetPlayListIDs() ([]string, error) {
	outInfo, err := wc.crawlerInfo()
	var arr []string
	err = codec.API.JsonAPI.Unmarshal(outInfo, &arr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	strSet := set.NewStringSet(len(arr) << 1)
	for _, s := range arr {
		split := strings.Split(s, "=")
		if len(split) != 2 {
			continue
		}
		strSet.Add(split[1])
	}
	return strSet.Iterator(), nil
}

func (wc *WangYiYunCrawler) getWangYiYunResult(id string) (*WangYiYunResult, error) {
	url := fmt.Sprintf(playlistBaseUrlFmt, id)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Error(closeErr)
		}
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result WangYiYunResult
	err = codec.API.JsonAPI.Unmarshal(bytes, &result)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &result, nil
}

func (wc *WangYiYunCrawler) ConversionToMusicWithPlaylistID(id string) ([]model.Music, error) {
	res, err := wc.getWangYiYunResult(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return wc.wangYiYunResultToMusic(res), nil
}

func (wc *WangYiYunCrawler) wangYiYunResultToMusic(refRes *WangYiYunResult) []model.Music {
	if refRes == nil {
		return nil
	}
	ret := make([]model.Music, 0, len(refRes.Result.Tracks))
	for _, track := range refRes.Result.Tracks {
		ret = append(ret, model.Music{
			ID:          track.ID,
			Name:        track.Name,
			Status:      0,
			Title:       track.Album.Name,
			CreatorID:   track.Artists[0].ID,
			CreatorName: track.Artists[0].Name,
			PlayTime:    track.Duration,
			ImageUrl:    track.Album.BlurPicUrl,
			PublishTime: fmt.Sprintf("%d", track.Album.PublishTime),
			UpdateTime:  tm.GetNowDateStr(),
		})
	}
	return ret
}