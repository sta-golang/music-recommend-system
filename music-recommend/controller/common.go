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

	creatorDetailUrl = "/creator/detail"
)

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
		timing := tm.FuncTiming(func() {
			fn(ctx)
		})
		log.Debugf("controller : %s timing : %v 毫秒", controllerName, timing.Milliseconds())
	}
}
