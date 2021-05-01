package algstat

import "github.com/sta-golang/music-recommend/service"

type userMusic struct {
}

func (us *userMusic) RegisterStatistics() {
	service.PubStatisticsService.Register("userMusic", &service.StatisticsFunc{
		ParseFunc:   us.parseStatistics,
		ProcessFunc: us.processStatistics,
	})
}

func (us *userMusic) processStatistics() {
}
func (us *userMusic) parseStatistics(bytes []byte) {
}
