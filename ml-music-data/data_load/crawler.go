package data_load

import (
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/cmd"
	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/music-recommend/model"
	"strings"
)

const (
	pythonCmd = "python"
	crawlerPythonPath = "./music_download/crawler.py"
)

type WangYiYunCrawler struct {
}

func (wc *WangYiYunCrawler) crawlerInfo() ([]byte, error) {
	command, err := cmd.ExecCmd(pythonCmd, crawlerPythonPath)
	if err != nil {
		return nil, err
	}
	message := command.OutMessage
	if len(message) > 0 {
		return message, nil
	}
	return nil, command.RunErr
}

func (wc *WangYiYunCrawler) GetPlayListIDs() ([]string, error) {
	outInfo, err := wc.crawlerInfo()
	var arr []string
	err = codec.API.JsonAPI.Unmarshal(outInfo, &arr)
	if err != nil {
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

func (wc *WangYiYunCrawler) getWangYiYunResult() {

}

func (wc *WangYiYunCrawler) ConversionToMusicWithPlaylistID(id string) ([]model.Music, error) {

}