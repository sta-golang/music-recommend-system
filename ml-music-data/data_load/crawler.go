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
	"time"
)

const (
	pythonCmd                = "python"
	pyCmd                    = "py"
	crawlerMusicPythonPath   = "./music_crawler/music_crawler.py"
	crawlerCreatorPythonPath = "./music_crawler/creator.py"
	playlistBaseUrlFmt       = "https://music.163.com/api/playlist/detail?id=%s"
	creatorResultFile        = "result.txt"
)

type MusicRecommendData struct {
	Musics   []model.Music
	Creators []model.Creator
}

type WangYiYunCrawler struct {
	creatorIDSet *set.StringSet
}

func (wc *WangYiYunCrawler) crawlerInfo() ([]byte, error) {
	command, err := cmd.ExecCmd(pyCmd, crawlerMusicPythonPath)
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

func (wc *WangYiYunCrawler) ConversionToDataWithPlaylistID(id string) (*MusicRecommendData, error) {
	ret := &MusicRecommendData{}
	res, err := wc.getWangYiYunResult(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 作者出了点问题
	creators, err := wc.wangYiYunResultToCreator(res)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ret.Creators = creators
	ret.Musics = wc.wangYiYunResultToMusic(res)
	return ret, nil
}

func (wc *WangYiYunCrawler) wangYiYunResultToCreator(refRes *WangYiYunResult) ([]model.Creator, error) {
	if refRes == nil {
		return nil, nil
	}
	//if wc.creatorIDSet == nil {
	//	wc.creatorIDSet = set.NewStringSet(3000)
	//}
	//ret := make([]model.Creator, 0, len(refRes.Result.Tracks))
	//for _, track := range refRes.Result.Tracks {
	//
	//	if wc.creatorIDSet.Contains(fmt.Sprintf("%d", track.Artists[0].ID)) {
	//		continue
	//	}
	//	command, err := cmd.ExecCmd(pyCmd, crawlerCreatorPythonPath, fmt.Sprintf("%d", track.Artists[0].ID))
	//	if err != nil {
	//		return nil, err
	//	}
	//	if command.RunErr != nil {
	//		log.Warn(command.RunErr)
	//	}
	//	// 不知道这里为什么控制台输出一直有问题
	//	//if len(command.OutMessage) <= 0 {
	//	//	continue
	//	//}
	//	fileOut, err := ioutil.ReadFile(creatorResultFile)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//	var creatorCrawler CrawlerCreator
	//	err = codec.API.JsonAPI.Unmarshal(fileOut, &creatorCrawler)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//	fmt.Println(string(fileOut))
	//
	//	similarCreatorStr := ""
	//	if len(creatorCrawler.SimilarCreator) > 0 {
	//		bytes, err := codec.API.JsonAPI.Marshal(creatorCrawler.SimilarCreator)
	//		if err != nil {
	//			log.Error(err)
	//			return nil, err
	//		}
	//		similarCreatorStr = str.BytesToString(bytes)
	//	}
	//
	//	ty := model.MusicianType
	//	if creatorCrawler.Superstar {
	//		ty = model.SuperstarType
	//	}
	//	ret = append(ret, model.Creator{
	//		ID:             track.Artists[0].ID,
	//		Name:           track.Artists[0].Name,
	//		Status:         0,
	//		ImageUrl:       creatorCrawler.ImageUrl,
	//		Description:    creatorCrawler.Description,
	//		SimilarCreator: similarCreatorStr,
	//		FansNum:        creatorCrawler.FansNum,
	//		Type:           ty,
	//		UpdateTime:     tm.GetNowDateTimeStr(),
	//	})
	//	fmt.Println(creatorCrawler.Description)
	//	wc.creatorIDSet.Add(fmt.Sprintf("%d", track.Artists[0].ID))
	//}
	ret := make([]model.Creator, 0)
	if wc.creatorIDSet == nil {
		wc.creatorIDSet = set.NewStringSet(50000)
	}
	for _, track := range refRes.Result.Tracks {
		creatorKey := fmt.Sprintf("%v-%v", track.Artists[0].ID, track.Artists[0].Name)
		wc.creatorIDSet.Add(creatorKey)
	}
	return ret, nil
}

func (wc WangYiYunCrawler) GetCreatorKeys() []string {
	if wc.creatorIDSet == nil {
		return nil
	}
	return wc.creatorIDSet.Iterator()
}

func (wc *WangYiYunCrawler) wangYiYunResultToMusic(refRes *WangYiYunResult) []model.Music {
	if refRes == nil {
		return nil
	}

	ret := make([]model.Music, 0, len(refRes.Result.Tracks))
	for _, track := range refRes.Result.Tracks {

		ret = append(ret, model.Music{
			ID:     track.ID,
			Name:   track.Name,
			Status: 0,
			Title:  track.Album.Name,

			CreatorID:   track.Artists[0].ID,
			TagNames:    strings.Join(refRes.Result.Tags, model.TagDelimiter),
			PlayTime:    track.Duration,
			ImageUrl:    track.Album.BlurPicUrl,
			PublishTime: tm.ParseDataTimeToStr(time.Unix(int64(track.Album.PublishTime)/1000, 0)),
			UpdateTime:  tm.GetNowDateStr(),
		})
	}
	return ret
}
