package data_load

import (
	"context"
	"errors"
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/async/dag"
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
	creatorMusicLimit        = 50
	creatorLimit             = 30
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
	ret.Musics = wc.wangYiYunResultToMusic(res, creators)
	return ret, nil
}

func (wc *WangYiYunCrawler) wangYiYunResultToCreator(refRes *WangYiYunResult) ([]model.Creator, error) {
	if refRes == nil {
		return nil, nil
	}

	ret := make([]model.Creator, 0, len(refRes.Result.Tracks))
	for _, track := range refRes.Result.Tracks {
		var endApiResult *APICreatorResult
		for _, at := range track.Artists {
			var apiResult APICreatorResult
			url := fmt.Sprintf("%s%s", APICrawlerName, apiResult.GetUrl(at.ID))
			bys, err := HttpGetFunc(url)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			err = codec.API.JsonAPI.Unmarshal(bys, &apiResult)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			if apiResult.Code != 200 {
				log.Error(fmt.Sprintf("Creator ID : %d, Code : %d Message : %s", at.ID,
					apiResult.Code, apiResult.Message))
				continue
			}
			if len(apiResult.Data.Artist.IdentifyTag) <= 0 || (len(apiResult.Data.Artist.IdentifyTag) > 0 &&
				apiResult.Data.Artist.IdentifyTag[0] != WangYiMusic) {
				endApiResult = &apiResult
				break
			}
			if endApiResult == nil {
				endApiResult = &apiResult
			}
		}
		var apiSimilarRes APISimilarResult
		if endApiResult == nil {
			log.Error("end nil ", endApiResult)
			return nil, errors.New("nil api Res")
		}
		bys, err := HttpGetFunc(fmt.Sprintf("%s%s", APICrawlerName, apiSimilarRes.GetUrl(endApiResult.Data.Artist.ID)))
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = codec.API.JsonAPI.Unmarshal(bys, &apiSimilarRes)
		if err != nil {
			return nil, err
		}
		similarStr := make([]string, 0, len(apiSimilarRes.Artists))
		for _, similar := range apiSimilarRes.Artists {
			similarStr = append(similarStr, fmt.Sprintf("%d", similar.ID))
		}
		creatorTy := model.CreatorUnknownType
		if endApiResult.Data.Artist.IdentifyTag == nil || (len(endApiResult.Data.Artist.IdentifyTag) > 0 &&
			endApiResult.Data.Artist.IdentifyTag[0] != WangYiMusic) {
			creatorTy = model.CreatorSuperstarType
		}
		if len(endApiResult.Data.Artist.IdentifyTag) > 0 &&
			endApiResult.Data.Artist.IdentifyTag[0] == WangYiMusic {
			creatorTy = model.CreatorMusicianType
		}
		if creatorTy == model.CreatorUnknownType {
			log.Warn("unknownType : ", endApiResult.Data.Artist.IdentifyTag)
		}
		ret = append(ret, model.Creator{
			ID:             endApiResult.Data.Artist.ID,
			Name:           endApiResult.Data.Artist.Name,
			Status:         0,
			ImageUrl:       endApiResult.Data.Artist.Cover,
			Description:    endApiResult.Data.Artist.BriefDesc,
			SimilarCreator: strings.Join(similarStr, model.TagDelimiter),
			FansNum:        0,
			Type:           creatorTy,
			UpdateTime:     "",
		})
	}
	return ret, nil
}

func (wc WangYiYunCrawler) GetCreatorKeys() []string {
	if wc.creatorIDSet == nil {
		return nil
	}
	return wc.creatorIDSet.Iterator()
}

func (wc *WangYiYunCrawler) wangYiYunResultToMusic(refRes *WangYiYunResult, creators []model.Creator) []model.Music {
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

			CreatorIDs:  "",
			TagNames:    strings.Join(refRes.Result.Tags, model.TagDelimiter),
			PlayTime:    track.Duration,
			ImageUrl:    track.Album.BlurPicUrl,
			PublishTime: tm.ParseDataTimeToStr(time.Unix(int64(track.Album.PublishTime)/1000, 0)),
			UpdateTime:  tm.GetNowDateStr(),
		})
	}
	return ret
}

func getCreatorListUrl(offset int) string {
	return fmt.Sprintf("/artist/list?type=-1&area=7&offset=%d&limit=%d", offset, creatorLimit)
}

func (wc *WangYiYunCrawler) CrawlerAllCreatorList() ([]model.Creator, error) {
	offset := 0
	creatorIDSet := set.NewStringSet(2000)
	creatorIDs := make([]int, 0, 1000)
	for {
		url := fmt.Sprintf("%s%s", APICrawlerName, getCreatorListUrl(offset))
		bys, err := HttpGetFunc(url)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var creatorIDList APICreatorListResult
		err = codec.API.JsonAPI.Unmarshal(bys, &creatorIDList)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, creatorData := range creatorIDList.Artists {
			if creatorIDSet.Contains(fmt.Sprintf("%d", creatorData.ID)) {
				continue
			}
			creatorIDs = append(creatorIDs, creatorData.ID)
			creatorIDSet.Add(fmt.Sprintf("%d", creatorData.ID))
		}
		if !creatorIDList.More {
			break
		}
		offset += creatorLimit
	}
	return wc.DoCrawlerAllCreatorList(creatorIDs)
}

func (wc *WangYiYunCrawler) CrawlerCreatorMusic(creatorID int) ([]APIMusicDetail, error) {
	var musicIDs []int
	offset := 0
	for {
		url := fmt.Sprintf("%s%s", APICrawlerName, getCreatorMusicUrl(creatorID, offset))
		bys, err := HttpGetFunc(url)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var musicResult APICreatorMusicResult
		err = codec.API.JsonAPI.Unmarshal(bys, &musicResult)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if musicIDs == nil {
			musicIDs = make([]int, 0, musicResult.Total)
		}
		for _, music := range musicResult.Songs {
			musicIDs = append(musicIDs, music.ID)
		}
		if !musicResult.More || offset > 10000 {
			break
		}
		offset += creatorMusicLimit
	}
	return wc.doCrawlerCreatorMusic(musicIDs)
}

func (wc *WangYiYunCrawler) doCrawlerCreatorMusic(musicIDs []int) ([]APIMusicDetail, error) {
	if len(musicIDs) <= 0 {
		return nil, nil
	}
	ret := make([]APIMusicDetail, 0, len(musicIDs))
	rootTask := dag.NewTask("root", func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {
		for _, musicID := range musicIDs {

			taskRet, err := helper.GetSubTaskRet(fmt.Sprintf("%d", musicID))
			if err != nil {
				log.Error(err)
				return nil, err
			}
			detail, ok := taskRet.([]APIMusicDetail)
			if !ok {
				log.Error("data error")
				return nil, errors.New("data error")
			}
			ret = append(ret, detail...)
		}
		return nil, nil
	})

	for _, musicID := range musicIDs {
		cruMusicID := musicID
		rootTask.AddSubTask(dag.NewTask(fmt.Sprintf("%d", cruMusicID), func(ctx context.Context, helper dag.TaskHelper) (interface{}, error) {

			url := fmt.Sprintf("%s%s", APICrawlerName, getMusicDetailUrl(cruMusicID))
			bys, err := HttpGetFunc(url)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			var res APIMusicDetailResult
			err = codec.API.JsonAPI.Unmarshal(bys, &res)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			if len(res.Songs) <= 0 {
				log.Error("songs len < 0")
				return nil, errors.New("songs len < 0")
			}
			return res.Songs, nil
		}))
	}
	dagTask := dag.NewDag(rootTask)
	if dagTask.Do(context.Background(), true) {
		log.Error("dag hasDependence")
		return nil, errors.New("dag hasDependence")
	}
	if _, rErr := rootTask.GetRet(); rErr != nil {
		return nil, rErr
	}
	return ret, nil
}

func getMusicDetailUrl(musicID int) string {
	return fmt.Sprintf("/song/detail?ids=%d", musicID)
}

func getCreatorMusicUrl(creatorID, offset int) string {
	return fmt.Sprintf("/artist/songs?id=%d&offset=%d&limit=%d", creatorID, offset, creatorMusicLimit)
}

func (wc *WangYiYunCrawler) DoCrawlerAllCreatorList(creatorIDs []int) ([]model.Creator, error) {
	if len(creatorIDs) <= 0 {
		return nil, nil
	}
	ret := make([]model.Creator, 0, len(creatorIDs))

	for _, creatorID := range creatorIDs {
		var apiResult APICreatorResult
		url := fmt.Sprintf("%s%s", APICrawlerName, apiResult.GetUrl(creatorID))
		bys, err := HttpGetFunc(url)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = codec.API.JsonAPI.Unmarshal(bys, &apiResult)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if apiResult.Code != 200 {
			log.Error(fmt.Sprintf("Creator ID : %d, Code : %d Message : %s", creatorID,
				apiResult.Code, apiResult.Message))
			continue
		}
		var apiSimilarRes APISimilarResult
		bys, err = HttpGetFunc(fmt.Sprintf("%s%s", APICrawlerName, apiSimilarRes.GetUrl(apiResult.Data.Artist.ID)))
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = codec.API.JsonAPI.Unmarshal(bys, &apiSimilarRes)
		if err != nil {
			return nil, err
		}
		similarStr := set.NewStringSet(len(apiSimilarRes.Artists))
		for _, similar := range apiSimilarRes.Artists {
			if similarStr.Contains(fmt.Sprintf("%d", similar.ID)) || similar.ID == 0 {
				continue
			}
			similarStr.Add(fmt.Sprintf("%d", similar.ID))
		}
		creatorTy := model.CreatorUnknownType
		if apiResult.Data.Artist.IdentifyTag == nil || (len(apiResult.Data.Artist.IdentifyTag) > 0 &&
			apiResult.Data.Artist.IdentifyTag[0] != WangYiMusic) {
			creatorTy = model.CreatorSuperstarType
		}
		if len(apiResult.Data.Artist.IdentifyTag) > 0 &&
			apiResult.Data.Artist.IdentifyTag[0] == WangYiMusic {
			creatorTy = model.CreatorMusicianType
		}
		if creatorTy == model.CreatorUnknownType {
			log.Warn("unknownType : ", apiResult.Data.Artist.IdentifyTag)
		}
		ret = append(ret, model.Creator{
			ID:             apiResult.Data.Artist.ID,
			Name:           apiResult.Data.Artist.Name,
			Status:         0,
			ImageUrl:       apiResult.Data.Artist.Cover,
			Description:    apiResult.Data.Artist.BriefDesc,
			SimilarCreator: strings.Join(similarStr.Iterator(), model.TagDelimiter),
			FansNum:        0,
			Type:           creatorTy,
			UpdateTime:     "",
		})
	}
	return ret, nil
}

func getPlaylistDetail(playlistID string) string {
	return fmt.Sprintf("/playlist/detail?id=%s", playlistID)
}

func (wc *WangYiYunCrawler) CrawlerPlaylistsDetail(playlistID string) (*APIPlaylistDetailResult, error) {
	var ret APIPlaylistDetailResult
	url := fmt.Sprintf("%s%s", APICrawlerName, getPlaylistDetail(playlistID))
	bys, err := HttpGetFunc(url)
	if err != nil {
		log.ConsoleLogger.Error(err)
		return nil, err
	}
	err = codec.API.JsonAPI.Unmarshal(bys, &ret)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &ret, nil
}
