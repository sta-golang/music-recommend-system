package common

import "errors"

const (
	// 第一位为xxx 二三位为xxxx 四五位为xxxx
	IdentifyingCurrencyErr = 50000
	IdentifyingCheckErr    = 50001
	IdentifyingSendErr     = 50002
	UserUnExistentErr      = 50003
	UserCheckErr           = 50004
	ServerCodecErr         = 50005

	DBFindErr   = 30000
	DBFindIDErr = 30001
	DBCreateErr = 30100
	DBUpdateErr = 30200
	DBDeleteErr = 30300

	ServiceUserRegistryErr = 20003

	Success         = 0
	NotFound        = 404
	NotFoundMessage = "资源未找到"

	//
	CurrencyErr = 99999
)

var (
	UserEmailExistsErr = errors.New("该邮箱已经被注册")
	ServerErr          = errors.New("服务器错误")
	CheckCodeErr       = errors.New("验证码错误")
	Affected           = errors.New("the number of affected rows is 0")
	InputEmailErr      = errors.New("账号输入错误或不存在")
	InputPasswordErr   = errors.New("密码输入错误")
	UserNotExistErr    = errors.New("用户不存在")
)
