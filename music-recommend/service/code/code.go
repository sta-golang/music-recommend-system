package code

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/service/cache"
	"github.com/sta-golang/music-recommend/service/email"
	"math/rand"
)

const (
	removeTime = 60 * 10
	codeFmt    = "gCode-%s"
)

var NewEmailIdentifyingService = func() Identifying {
	return onceEmailIdentify
}

var onceEmailIdentify = &EmailIdentify{}

type EmailIdentify struct {
}

func (ei *EmailIdentify) Generate() string {
	return fmt.Sprintf("%06v", rand.Int31n(1000000))
}

func (ei *EmailIdentify) SendCode(key, code string) error {
	if err := ei.doSendCode(key, code); err != nil {
		log.Error(err)
		return err
	}
	cache.PubCacheService.Set(fmt.Sprintf(codeFmt, key), code, removeTime, cache.Ten)
	return nil
}

func (ei *EmailIdentify) doSendCode(key, code string) error {
	// 主题  用户名 类别 验证码 类别 类别
	body := fmt.Sprintf(emailHtmlFmt, "验证码", key, common.PublicName, code, common.PublicName, common.PublicName)
	err := email.PubEmailService.SendEmail(codeSubject, body, key)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (ei *EmailIdentify) CheckCode(key, code string) (bool, error) {
	ret, ok := cache.PubCacheService.Get(fmt.Sprintf(codeFmt, key))
	if !ok {
		return false, nil
	}
	return ret == code, nil
}
