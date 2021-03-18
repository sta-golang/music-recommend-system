package service

import "github.com/sta-golang/go-lib-utils/pool/workerpool"

func init()  {
	err := workerpool.Submit(PubStatisticsService.Run)
	if err != nil {
		go PubStatisticsService.Run()
	}
}
