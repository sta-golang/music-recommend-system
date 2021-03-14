package data_load

import (
	"github.com/sta-golang/go-lib-utils/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	reTryNum = 300
)

func HttpGetFunc(url string) ([]byte, error) {
	var retErr error
	for i := 0; i < reTryNum; i++ {
		if retErr != nil {
			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		}
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err)
			if retErr == nil {
				retErr = err
			}
			continue
		}
		defer resp.Body.Close()
		retBys, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			if retErr == nil {
				retErr = err
			}
			continue
		}
		return retBys, nil
	}

	return nil, retErr
}
