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

func (pc *playlistController) TestPlaylistAPI(ctx *fasthttp.RequestCtx) {
	api := string(ctx.FormValue("api"))
	username := "554285007@qq.com"
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	playlistID := ctx.QueryArgs().GetUintOrZero("id")
	if api == "add" {
		name := string(ctx.FormValue("name"))
		err := service.PubPlaylistService.AddPlaylistForUserWithCache(reqCtx, name, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
		return
	} else if api == "delete" {
		err := service.PubPlaylistService.DeletePlaylistForUserWithCache(reqCtx, playlistID, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
	} else if api == "get" {
		data, err := service.PubPlaylistService.GetPlaylistForUserWithCache(reqCtx, username)
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

func (pc *playlistController) TestPlaylistMusicAPI(ctx *fasthttp.RequestCtx) {
	api := string(ctx.FormValue("api"))
	username := "554285007@qq.com"
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	playlistID := ctx.QueryArgs().GetUintOrZero("pid")
	musicID := ctx.QueryArgs().GetUintOrZero("mid")
	if api == "add" {
		err := service.PubPlaylistService.AddMusicToPlaylist(reqCtx, musicID, playlistID, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
		return
	} else if api == "delete" {
		err := service.PubPlaylistService.DeleteMusicForPlaylist(reqCtx, musicID, playlistID, username)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
	} else if api == "get" {
		data, err := service.PubPlaylistService.GetPlaylistMusicWithCache(reqCtx, playlistID)
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
	musics, _ := service.PubPlaylistService.GetPlaylistMusicWithCache(ctx, id)
	retData := dto.UserAndPlaylist{
		Playlist: playlist,
		User:     user,
		Musics:   musics,
	}
	WriterResp(ctx, NewRetData(successCode, success, retData).ToJson())
}

var noneUser = &model.User{
	Name: "该用户已注销",
}

func (pc *playlistController) GetUserPlaylist(ctx *fasthttp.RequestCtx) {
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	username := str.BytesToString(ctx.FormValue("username"))
	playlists, err := service.PubPlaylistService.GetPlaylistForUser(reqCtx, username)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(common.ServerCodecErr, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, playlists).ToJson())
}

func (pc *playlistController) GetPlaylistHot(ctx *fasthttp.RequestCtx) {
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	playlists, err := service.PubPlaylistService.GetHotPlaylistWithCache(reqCtx)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(common.ServerCodecErr, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, playlists).ToJson())
}

func (pc *playlistController) GetPlaylistMusic(ctx *fasthttp.RequestCtx) {
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	id := ctx.QueryArgs().GetUintOrZero("id")
	if id == 0 {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	retData, err := service.PubPlaylistService.GetPlaylistMusicWithCache(reqCtx, id)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, err.Error()).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, retData).ToJson())
}

func (pc *playlistController) AddPlaylistForUser(ctx *fasthttp.RequestCtx) {
	user, ok := haveAuthority(ctx)
	if !ok {
		return
	}
	if user == nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, "用户异常").ToJson())
		return
	}
	name := str.BytesToString(ctx.FormValue("name"))
	if name == "" {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, "歌单名不能为空").ToJson())
		return
	}
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	err := service.PubPlaylistService.AddPlaylistForUserWithCache(reqCtx, name, user.Username)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, err.Error()).ToJson())
		return
	}
	WriterResp(ctx, NewRetDataForErrorAndMessage(successCode, success).ToJson())
}
