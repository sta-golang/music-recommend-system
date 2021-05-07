package controller

import (
	"context"
	"net/http"

	"github.com/sta-golang/music-recommend/service"
	"github.com/valyala/fasthttp"
)

type musicController struct {
}

var onceMusicController = musicController{}

func NewMusicController() *musicController {
	return &onceMusicController
}

// CreatorMusics 作者的相关音乐
func (mc *musicController) CreatorMusics(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	id := args.GetUintOrZero("id")
	if id == 0 {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	page := args.GetUintOrZero("page")
	creator, sErr := service.PubMusicService.GetMusicForCreator(id, page)
	if sErr != nil && sErr.Err != nil {
		WriterResp(ctx, NewRetDataForErrAndMessage(sErr, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, creator).ToJson())
}

func (mc *musicController) GetAllMusics(ctx *fasthttp.RequestCtx) {
	page := ctx.QueryArgs().GetUintOrZero("page")
	musics, sErr := service.PubMusicService.GetAllMusicWithCache(context.Background(), page)
	if sErr != nil && sErr.Err != nil {
		WriterResp(ctx, NewRetDataForErrAndMessage(sErr, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, musics).ToJson())
}

// GetMusic 获取音乐
func (mc *musicController) GetMusic(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	id := args.GetUintOrZero("id")
	if id == 0 {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	music, sErr := service.PubMusicService.GetMusic(id)
	if sErr != nil && sErr.Err != nil {
		WriterResp(ctx, NewRetDataForErrAndMessage(sErr, serverSelectErrMessage).ToJson())
		return
	}
	user := getSessionUser(ctx)
	service.PubUserMusicService.StatMusicForUser(reqCtx, user, id)
	WriterResp(ctx, NewRetData(successCode, success, music).ToJson())
}
