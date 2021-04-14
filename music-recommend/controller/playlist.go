package controller

import (
	"net/http"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/controller/dto"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
	"github.com/valyala/fasthttp"
)

type playlistController struct {
}

var oncePlaylistController = playlistController{}

func NewPlaylistController() *playlistController {
	return &oncePlaylistController
}

func (pc *playlistController) TestAPI(ctx *fasthttp.RequestCtx) {
	api := string(ctx.FormValue("api"))
	username := "554285007@qq.com"
	ReqCtx := RequestContext(ctx)
	defer DestroyContext(ReqCtx)
	playlistID := ctx.QueryArgs().GetUintOrZero("id")
	if api == "add" {
		name := string(ctx.FormValue("name"))
		err := service.PubPlaylistService.AddPlaylistForUserWithCache(ReqCtx, name, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
		return
	} else if api == "delete" {
		err := service.PubPlaylistService.DeletePlaylistForUserWithCache(ReqCtx, playlistID, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
	} else if api == "get" {
		data, err := service.PubPlaylistService.GetPlaylistForUserWithCache(ReqCtx, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		bytes, _ := codec.API.JsonAPI.Marshal(data)
		ctx.Write(bytes)
	} else {
		ctx.WriteString("404 not page")
	}
}

func (pc *playlistController) GetPlaylistDetail(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().GetUintOrZero("id")
	if id == 0 {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	playlist, err := service.PubPlaylistService.GetPlaylistDetailWithCache(reqCtx, id)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, err.Error()).ToJson())
		return
	}
	if playlist == nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusNotFound, common.NotFoundMessage).ToJson())
	}
	user, _ := service.PubUserService.QueryUserWithCache(playlist.Username)
	if user == nil {
		user = noneUser
	}
	retData := dto.UserAndPlaylist{
		Playlist: playlist,
		User:     user,
	}
	WriterResp(ctx, NewRetData(successCode, success, retData).ToJson())
}

var noneUser = &model.User{
	Name: "该用户已注销",
}

func (pc *playlistController) GetUserPlaylist(ctx *fasthttp.RequestCtx) {
	ReqCtx := RequestContext(ctx)
	defer DestroyContext(ReqCtx)
	username := str.BytesToString(ctx.FormValue("username"))
	playlists, err := service.PubPlaylistService.GetPlaylistForUser(ReqCtx, username)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(common.ServerCodecErr, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, playlists).ToJson())
}
