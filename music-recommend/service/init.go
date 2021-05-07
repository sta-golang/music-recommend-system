package service

import "github.com/sta-golang/go-lib-utils/pool/workerpool"

func init() {
	PubMusicService.RegisterStatistics()
	PubUserService.RegisterStatistics()
	PubUserMusicService.RegisterStatistics()
	err := workerpool.Submit(PubStatisticsService.Run)
	if err != nil {
		go PubStatisticsService.Run()
	}
}
