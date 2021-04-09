package controller

import (
	"github.com/sta-golang/go-lib-utils/codec"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/valyala/fasthttp"
)

const (
	paramsErrMessage       = "参数错误"
	serverSelectErrMessage = "服务器查询出错 请联系管理员QQ:63237777"
	success                = "success"
	successCode            = 0
	postDataErrMessage     = "提交数据表单异常"
	forbiddenErrMessage    = "没有权限"
	tokenTimeOutErrMessage = "凭证已过期请重新登录"
	waitMessage            = "请稍后再重试"
	sendCodeMessage        = "发送验证码成功！\n如果没有收到请查看垃圾邮件"
	tokenStr               = "sta-token"

	creatorDetailUrl = "/creator/detail"
	creatorList      = "/creator/list"
	musicDetails     = "/music/details"
	creatorMusic     = "/creator/music"
	userRegister     = "/user/register"
	userLogin        = "/user/login"
	userCode         = "/user/code"
	userInfo         = "/user/me"
	musicAll         = "/music/all"
)

/**
这里全是一些视图层需要返回的数据定义。方便一些
*/

var DebugLevelName = log.GetLevelName(log.DEBUG)

type RetData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewRetData(code int, message string, data interface{}) *RetData {
	return &RetData{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewRetDataForErr(err *er.Error) *RetData {
	return &RetData{
		Code:    err.Code,
		Message: err.Err.Error(),
		Data:    nil,
	}
}

func NewRetDataForErrorAndMessage(code int, message string) *RetData {
	return &RetData{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func NewRetDataForErrAndMessage(err *er.Error, message string) *RetData {
	return &RetData{
		Code:    err.Code,
		Message: message,
		Data:    nil,
	}
}

func (rd *RetData) ToJson() []byte {
	bytes, _ := codec.API.JsonAPI.Marshal(rd)
	return bytes
}

func WriterResp(ctx *fasthttp.RequestCtx, bys []byte) {
	_, err := ctx.Write(bys)
	if err != nil {
		log.Error(err)
	}

}

func TimeController(controllerName string, fn fasthttp.RequestHandler) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		if log.GetLevel() == DebugLevelName {
			timing := tm.FuncTiming(func() {
				fn(ctx)
			})
			log.Debugf("controller : %s timing : %v 毫秒", controllerName, timing.Milliseconds())
		} else {
			fn(ctx)
		}
	}
}

func CORSHandler(fn fasthttp.RequestHandler) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		canCORS(ctx)
		fn(ctx)
	}
}

func canCORS(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "PUT,GET,POST,DELETE,OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Set("content-type", "application/json")
}
