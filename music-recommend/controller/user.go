package controller

import (
	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/controller/dto"
	"github.com/sta-golang/music-recommend/service"
	"github.com/sta-golang/music-recommend/service/cache"
	"github.com/sta-golang/music-recommend/service/verify"
	"github.com/valyala/fasthttp"
	"net/http"
)

const (
	codeUser = "code_%s_u"
)

type userController struct {
}

var onceUserController = userController{}

func NewUserController() *userController {
	return &onceUserController
}

func (uc *userController) MeInfo(ctx *fasthttp.RequestCtx) {
	token := str.BytesToString(ctx.Request.Header.Peek(tokenStr))
	if token == "" {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, forbiddenErrMessage).ToJson())
		return
	}
	auth, ok, err := verify.NewJWTService().VerifyAuth(token)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, err.Error()).ToJson())
		return
	}
	if !ok {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, tokenTimeOutErrMessage).ToJson())
		return
	}
	info, exist := service.PubUserService.MeInfo(auth)
	if !exist {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusForbidden, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, info).ToJson())
}

func (uc *userController) SendCode(ctx *fasthttp.RequestCtx) {
	canCORS(ctx)
	username := str.BytesToString(ctx.FormValue("username"))
	if _, ok := cache.PubCacheService.Get(username); ok {
		WriterResp(ctx, NewRetData(successCode, waitMessage, nil).ToJson())
		return
	}
	if sErr := service.PubUserService.SendCodeForUser(username); sErr != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusInternalServerError, sErr.String()).ToJson())
		return
	}
	cache.PubCacheService.Set(username, true, 60, cache.Two)
	WriterResp(ctx, NewRetData(successCode, sendCodeMessage, nil).ToJson())
}

func (uc *userController) Login(ctx *fasthttp.RequestCtx) {
	canCORS(ctx)
	username := str.BytesToString(ctx.FormValue("username"))
	password := str.BytesToString(ctx.FormValue("password"))
	readme := ctx.PostArgs().GetBool("readme")
	if username == "" || password == "" || len(username) > 25 || len(password) > 25 ||
		len(username) < 6 || len(password) < 6 {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	token, sErr := service.PubUserService.Login(username, password, readme)
	if sErr != nil {
		log.Error(sErr)
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, sErr.String()).ToJson())
		return
	}
	ctx.Response.Header.Set(tokenStr, token)
	WriterResp(ctx, NewRetData(successCode, success, nil).ToJson())
}

func (uc *userController) Register(ctx *fasthttp.RequestCtx) {
	canCORS(ctx)
	var rUser dto.RegisterUser
	err := codec.API.JsonAPI.Unmarshal(ctx.PostBody(), &rUser)
	if err != nil {
		log.Error(err)
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, postDataErrMessage).ToJson())
		return
	}
	if rUser.User.Username == "" || rUser.User.Password == "" || len(rUser.User.Username) > 25 || rUser.Code == "" ||
		len(rUser.User.Password) > 25 || len(rUser.User.Username) < 6 || len(rUser.User.Password) < 6 {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	sErr := service.PubUserService.Register(&rUser.User, rUser.Code)
	if sErr != nil {
		log.Error(sErr)
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, sErr.String()).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, nil).ToJson())
}
