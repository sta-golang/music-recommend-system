package controller

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/valyala/fasthttp"
)

func GlobalRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET(creatorDetailUrl, TimeController(creatorDetailUrl, NewCreatorController().GetCreator))
	return router
}

// index 页
func Index(ctx *fasthttp.RequestCtx) {
	_, err := fmt.Fprint(ctx, "Welcome STA-Golang Music-Recommend")
	if err != nil {
		log.Error(err)
	}
}
