package controller

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/valyala/fasthttp"
)

// GlobalRouter 返回 路由器
func GlobalRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", Index)

	//router.GET(creatorDetailUrl, TimeController(creatorDetailUrl, NewCreatorController().GetCreator))
	router.GET(creatorDetailUrl, TimeController(creatorDetailUrl, NewCreatorController().GetCreator))
	router.GET(creatorList, TimeController(creatorList, NewCreatorController().GetCreators))
	router.GET(musicDetails, TimeController(musicDetails, NewMusicController().GetMusic))
	router.GET(creatorMusic, TimeController(creatorMusic, NewMusicController().CreatorMusics))
	router.POST(userRegister, TimeController(userRegister, NewUserController().Register))
	router.POST(userLogin, TimeController(userLogin, NewUserController().Login))
	router.POST(userCode, TimeController(userCode, NewUserController().SendCode))
	router.GET(userInfo, TimeController(userInfo, NewUserController().MeInfo))
	router.GET(musicAll, TimeController(musicAll, NewMusicController().GetAllMusics))
	router.GET("/test/playlist", NewPlaylistController().TestPlaylistAPI)
	router.GET("/test/playlistMusic", NewPlaylistController().TestPlaylistMusicAPI)
	router.GET(playlistMusic, TimeController(playlistMusic, NewPlaylistController().GetPlaylistMusic))
	router.GET(playlistUser, TimeController(playlistUser, NewPlaylistController().GetUserPlaylist))
	router.GET(plsylistDetail, TimeController(plsylistDetail, NewPlaylistController().GetPlaylistDetail))
	return router
}

// index 页
func Index(ctx *fasthttp.RequestCtx) {
	_, err := fmt.Fprint(ctx, "Welcome STA-Golang Music-Recommend")
	if err != nil {
		log.Error(err)
	}
}
