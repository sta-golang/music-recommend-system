package controller

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/sta-golang/go-lib-utils/codec"
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"github.com/sta-golang/go-lib-utils/source"
	str "github.com/sta-golang/go-lib-utils/str"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
	"github.com/sta-golang/music-recommend/service/email"
	"github.com/sta-golang/music-recommend/service/verify"
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
	emailFormErr           = "邮箱格式错误"
	sessionIDStr           = "sessionID"

	playlistUser       = "/playlist/user"
	plsylistDetail     = "/playlist/detail"
	playlistMusic      = "/playlist/music"
	playlistAdd        = "/playlist/add"
	creatorDetailUrl   = "/creator/detail"
	creatorList        = "/creator/list"
	musicDetails       = "/music/details"
	creatorMusic       = "/creator/music"
	userRegister       = "/user/register"
	userLogin          = "/user/login"
	userCode           = "/user/code"
	userInfo           = "/user/me"
	musicAll           = "/music/all"
	recommendMusicList = "/recommend/list"
)

/**
这里全是一些视图层需要返回的数据定义。方便一些
*/

var DebugLevelName = log.GetLevelName(log.DEBUG)
var contextKeyMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string)
	},
}

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

func RequestContext(reqCtx *fasthttp.RequestCtx) context.Context {
	keyMap := contextKeyMapPool.Get().(map[string]string)
	keyMap["requestID"] = common.UUID()
	keyMap["scene"] = str.BytesToString(reqCtx.Path())
	ipBys, _ := reqCtx.RemoteIP().MarshalText()
	keyMap["ipAddr"] = str.BytesToString(ipBys)
	keyMap["user"] = ""
	token := str.BytesToString(reqCtx.Request.Header.Peek(tokenStr))
	if token != "" {
		auth, _, _ := verify.NewJWTService().VerifyAuth(token)
		if auth != "" {
			keyMap["user"] = auth
		}
	}
	return log.LogContextKeyMap(nil, keyMap)
}

func GetSessionID(reqCtx *fasthttp.RequestCtx) string {
	return str.BytesToString(reqCtx.Request.Header.Peek(sessionIDStr))
}

func RequestContextAndUser(reqCtx *fasthttp.RequestCtx, username string) context.Context {
	keyMap := contextKeyMapPool.Get().(map[string]string)
	keyMap["requestID"] = common.UUID()
	keyMap["scene"] = str.BytesToString(reqCtx.Path())
	ipBys, _ := reqCtx.RemoteIP().MarshalText()
	keyMap["ipAddr"] = str.BytesToString(ipBys)
	keyMap["user"] = username
	return log.LogContextKeyMap(nil, keyMap)
}

func DestroyContext(ctx context.Context) {
	if ctx == nil {
		return
	}
	if val := ctx.Value(log.LogContextKey); val != nil {
		contextKeyMapPool.Put(val)
	}
	ctx = nil
}

func getSessionUser(ctx *fasthttp.RequestCtx) *model.User {
	token := str.BytesToString(ctx.Request.Header.Peek(tokenStr))
	if token == "" {
		return nil
	}
	auth, ok, err := verify.NewJWTService().VerifyAuth(token)
	if err != nil {
		return nil
	}
	if !ok {
		return nil
	}
	info, exist := service.PubUserService.MeInfo(auth)
	if !exist {
		return nil
	}
	return info
}

func haveAuthority(ctx *fasthttp.RequestCtx) (*model.User, bool) {
	token := str.BytesToString(ctx.Request.Header.Peek(tokenStr))
	if token == "" {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, forbiddenErrMessage).ToJson())
		return nil, false
	}
	auth, ok, err := verify.NewJWTService().VerifyAuth(token)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, err.Error()).ToJson())
		return nil, false
	}
	if !ok {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, tokenTimeOutErrMessage).ToJson())
		return nil, false
	}
	info, exist := service.PubUserService.MeInfo(auth)
	if !exist {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, common.UserEmailNotLogin.Error()).ToJson())
		return nil, false
	}
	return info, true
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

func ServerHandler(fn fasthttp.RequestHandler) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if er := recover(); er != nil {
				ctx.WriteString("目前此功能可能出现部分问题，请留意后续版本或者联系管理员qq:63237777")
				source.Sync()
				debugInfo := str.BytesToString(debug.Stack())
				debugInfo = fmt.Sprintf("err %v \n %s", er, debugInfo)
				log.Fatal(debugInfo)
				fn := func() {
					subject := fmt.Sprintf(email.ServerMsg, "Fatal")
					email.PubEmailService.SendEmail(subject, debugInfo, config.GlobalConfig().EmailConfig.Email)
				}
				if err := workerpool.Submit(fn); err != nil {
					go fn()
				}
			}
		}()
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
