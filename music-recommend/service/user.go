package service

import (
	"fmt"
	"regexp"
	"time"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure"
	"github.com/sta-golang/go-lib-utils/codec"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service/cache"
	cd "github.com/sta-golang/music-recommend/service/code"
	"github.com/sta-golang/music-recommend/service/verify"
	"github.com/valyala/bytebufferpool"
)

const (
	defExpireTime       = time.Hour * 5
	defReadmeExpireTime = time.Hour * 24 * 30

	userCacheKeyFmt        = "user_%s_u"
	userSessionCacheKeyFmt = "user_%s_session"
)

type userService struct {
}

var PubUserService = &userService{}
var passwordRegexp = regexp.MustCompile("^[a-zA-Z]\\S{7,19}$")
var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var nameRegexp = regexp.MustCompile("^\\S{2,12}$")

func (us *userService) Register(user *model.User, code string) *er.Error {
	if user == nil {
		return er.NewError(common.UserCheckErr, fmt.Errorf("用户数据为空"))
	}
	if code == "" {
		return er.NewError(common.IdentifyingCheckErr, common.CodeNotExistErr)
	}
	if sErr := us.checkUser(user); sErr != nil {
		return sErr
	}
	var err error
	flag, err := cd.NewEmailIdentifyingService().CheckCode(user.Username, code)
	if err != nil {
		log.Error(err)
		return er.NewError(common.ServerCodecErr, common.CodeSendErr)
	}
	if !flag {
		return er.NewError(common.UserCheckErr, common.CheckCodeErr)
	}
	user.Password, err = codec.API.CryptoAPI.Encode(user.Password)
	if err != nil {
		log.Error(err)
		return er.NewError(common.ServerCodecErr, common.ServerErr)
	}
	err = model.NewUserMysql().Insert(user)
	if err != nil {
		log.Error(err)
		return er.NewError(common.DBCreateErr, common.ServerErr)
	}
	return nil
}

func (us *userService) SendCodeForUser(username string) *er.Error {
	if ok, sErr := us.checkUserExistAndSetCache(&username); ok {
		return sErr
	}
	fn := func() {
		err := cd.NewEmailIdentifyingService().SendCode(username, cd.NewEmailIdentifyingService().Generate())
		if err != nil {
			log.Error(err)
		}
	}
	if err := workerpool.Submit(fn); err != nil {
		go fn()
	}
	return nil
}

func (us *userService) MeInfo(username string) (*model.User, bool) {
	ret, ok := us.meInfo(username)
	if ret == nil {
		return nil, false
	}
	user := ret
	retUser := *user
	retUser.Password = ""
	return &retUser, ok
}

func (us *userService) meInfo(username string) (*model.User, bool) {
	ret, ok := cache.PubCacheService.Get(fmt.Sprintf(userSessionCacheKeyFmt, username))
	if ret == nil {
		return nil, false
	}
	return ret.(*model.User), ok
}

func (us *userService) checkUserExistAndSetCache(username *string) (bool, *er.Error) {
	user, sErr := us.QueryUserWithCache(*username)
	if sErr != nil {
		return true, sErr
	}
	if user != nil {
		return true, er.NewError(common.UserCheckErr, common.UserEmailExistsErr)
	}
	return false, nil
}

func (us *userService) QueryUserWithCache(username string) (*model.User, *er.Error) {
	key := fmt.Sprintf(userCacheKeyFmt, username)
	if val, ok := cache.PubCacheService.Get(key); ok && val != nil {
		return val.(*model.User), nil
	}
	ret, err := common.SingleRunGroup.Do(key, func() (interface{}, error) {
		queryObj, err := model.NewUserMysql().SelectUserForUserName(username)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if queryObj == nil {
			cache.PubCacheService.Set(key, nil, 60*60*72, cache.Five)
			return nil, nil
		}
		cache.PubCacheService.Set(key, queryObj, 60*60*72, cache.One)
		return queryObj, nil
	})
	if err != nil {
		return nil, er.NewError(common.DBFindErr, common.ServerErr)
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*model.User), nil
}

func (us *userService) Login(username, password string, readme bool) (token string, sErr *er.Error) {
	var user *model.User
	var err error
	if user, sErr = us.QueryUserWithCache(username); sErr != nil {
		return "", sErr
	}
	if user == nil {
		return "", er.NewError(common.UserUnExistentErr, common.UserNotExistErr)
	}
	if !codec.API.CryptoAPI.Check(password, user.Password) {
		return "", er.NewError(common.UserCheckErr, common.LoginUserErr)
	}
	expireTime := defExpireTime
	if readme {
		expireTime = defReadmeExpireTime
	}
	token, err = verify.NewJWTService().CreateToken(username, expireTime)
	if err != nil {
		log.Error(err)
		return "", nil
	}
	buff := bytebufferpool.Get()
	if _, err = buff.WriteString(username); err != nil {
		log.Error(err)
		bytebufferpool.Put(buff)
		return token, nil
	}
	PubStatisticsService.Statistics(us.GetName(), buff, false)
	cache.PubCacheService.Set(fmt.Sprintf(userSessionCacheKeyFmt, username), user, int(expireTime.Seconds()), cache.Ten)
	return token, nil
}

func (us *userService) IsEmailFmt(email string) bool {
	return emailRegexp.MatchString(email)
}

func (us *userService) checkUser(user *model.User) *er.Error {
	if !emailRegexp.MatchString(user.Username) {
		return er.NewError(common.UserCheckErr, common.InputEmailErr)
	}
	if !passwordRegexp.MatchString(user.Password) {
		return er.NewError(common.UserCheckErr, common.InputPasswordErr)
	}
	if !nameRegexp.MatchString(user.Name) {
		return er.NewError(common.UserCheckErr, common.InputUserNameErr)
	}
	return nil
}

func (us *userService) GetName() string {
	return "userService"
}

func (us *userService) RegisterStatistics() {
	queue := data_structure.NewQueue()
	PubStatisticsService.Register(us.GetName(), &StatisticsFunc{
		ParseFunc: func(bytes []byte) {
			queue.Push(str.BytesToString(bytes))
		},
		ProcessFunc: func() {
			for !queue.Empty() {
				username := queue.Pop().(string)
				check, err := model.NewUserMysql().ReSetStatistics(username)
				if err != nil {
					log.Error(err)
					continue
				}
				if check {
					continue
				}
				err = model.NewUserMysql().UserLogin(username)
				if err != nil {
					log.Error(err)
					continue
				}
			}
		},
	})
}
