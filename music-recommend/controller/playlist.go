package controller

import (
	"github.com/sta-golang/go-lib-utils/codec"
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
	userID := 2
	ReqCtx := RequestContext(ctx)
	defer DestroyContext(ReqCtx)
	playlistID := ctx.QueryArgs().GetUintOrZero("id")
	if api == "add" {
		name := string(ctx.FormValue("name"))
		err := service.PubPlaylistService.AddPlaylistForUserWithCache(ReqCtx, name, userID)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
		return
	} else if api == "delete" {
		err := service.PubPlaylistService.DeletePlaylistForUserWithCache(ReqCtx, playlistID, userID)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		ctx.WriteString("ok")
	} else if api == "get" {
		data, err := service.PubPlaylistService.GetPlaylistForUserWithCache(ReqCtx, userID)
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
