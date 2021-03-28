package benchmark_test

import (
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/service"
	"testing"
)

func BenchmarkService(b *testing.B)  {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, e := service.PubCreatorService.GetCreatorDetail(7214)
			log.Info(e)
		}
	})
}