package controller

import (
	"github.com/sta-golang/music-recommend/service"
	"github.com/valyala/fasthttp"
	"net/http"
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

// GetMusic 获取音乐
func (mc *musicController) GetMusic(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
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
	WriterResp(ctx, NewRetData(successCode, success, music).ToJson())
}
