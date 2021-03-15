package service

import "github.com/sta-golang/go-lib-utils/pool/workerpool"

func init()  {
	err := workerpool.Submit(PubMusicService.statisticsService)
	if err != nil {
		go PubMusicService.statisticsService()
	}
}
