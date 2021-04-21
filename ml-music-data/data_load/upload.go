package data_load

import (
	"context"
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/go-lib-utils/cmd"
	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/ml-music-data/utils"
	"github.com/sta-golang/music-recommend/model"
	"os"
	"strconv"
)

const (
	downloadPYPath = "./music_crawler/download.py"

	cosPath = "music/%d/%s-%d.mp3"
	localPath = "D:\\music\\%s-%d.mp3"
)

type UPLoad struct {
}

func (up *UPLoad) DownloadPy(musics []model.Music) (sucessMusic []model.Music, notHasMusicIDs []string, err error) {
	arr := make([]string, len(musics))
	musicSet := set.NewStringSet(len(musics))
	for i := 0; i < len(musics); i++ {
		musicSet.Add(strconv.Itoa(musics[i].ID))
		arr[i] = fmt.Sprintf("%d=%s", musics[i].ID, musics[i].Name)
	}
	bys, err := codec.API.JsonAPI.Marshal(arr)
	if err != nil {
		return
	}

	jsonStr := str.BytesToString(bys)
	command, err := cmd.ExecCmd(pyCmd, downloadPYPath,jsonStr)
	if err != nil {
		log.Error(err)
		return
	}
	message := command.OutMessage
	if command.RunErr != nil {
		log.Warn(command.RunErr)
	}
	if len(message) > 0 {
		var sucessIDs []string
		err = codec.API.JsonAPI.Unmarshal(message, &sucessIDs)
		if err != nil {
			return
		}
		musicSet.Remove(sucessIDs...)
		notHasMusicIDs = musicSet.Iterator()
		for k := range musics {
			for i := range sucessIDs {
				if sucessIDs[i] == strconv.Itoa(musics[k].ID) {
					sucessMusic = append(sucessMusic, musics[k])
					break
				}
			}
		}

	}
	return
}

func (up *UPLoad) UploadMusic(ctx context.Context, music *model.Music) error {
	cosPath := fmt.Sprintf(cosPath, music.ID % 100, music.Name,music.ID)
	LocalFileName := fmt.Sprintf(localPath, music.Name, music.ID)
	file, err := os.Open(LocalFileName)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		dErr := os.Remove(LocalFileName)
		if dErr != nil {
			log.Error(dErr)
		}
	}()
	err = utils.CosPutObject(ctx, cosPath, file, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	dbUrl := fmt.Sprintf("%s/%s", utils.CosURL(), cosPath)
	music.Status = model.MusicHasMusicUrlStatus
	music.MusicUrl = dbUrl
	_, err = model.NewMusicMysql().UpdateMusic(music)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
