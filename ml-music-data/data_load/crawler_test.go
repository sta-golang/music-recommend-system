package data_load

import (
	"fmt"
	"testing"
	"time"
)

func TestWangYiYunCrawler_ConversionToDataWithPlaylistID(t *testing.T) {
	wy := &WangYiYunCrawler{}
	result, err := wy.getWangYiYunResult("2557908184")
	fmt.Println(result, err)
	fmt.Println(result.Result.Tracks[0].Duration)
	fmt.Println(result.Result.Tracks[0].Album.PublishTime)
	fmt.Println(time.Unix(int64(result.Result.Tracks[0].Album.PublishTime)/1000, 0))
}
