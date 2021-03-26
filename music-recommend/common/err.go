package common

import "errors"

/**
业务错误的一些错误码定义。返回后方便定位错误位置
*/
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
	InputEmailErr      = errors.New("账号输入错误")
	LoginUserErr       = errors.New("账号密码错误")
	InputPasswordErr   = errors.New("密码输入错误")
	InputUserNameErr   = errors.New("昵称输入错误")
	UserNotExistErr    = errors.New("用户不存在")
	CodeNotExistErr    = errors.New("验证码为空")
	CodeSendErr        = errors.New("系统验证码发送出错,请联系管理员QQ:6323777")
	CreatorNotExistErr = errors.New("作者不存在")
	UnFollowRepeatWarn = errors.New("请勿重复取消关注")
	FollowRepeatWarn = errors.New("请勿重复关注")
)
