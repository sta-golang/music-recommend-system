package controller

import (
	"fmt"

	"github.com/sta-golang/music-recommend/common"
	"github.com/sta-golang/music-recommend/controller/dto"
	"github.com/sta-golang/music-recommend/feed"
	"github.com/sta-golang/music-recommend/model"
	"github.com/valyala/fasthttp"
)

type recommendController struct {
}

var onceRecommendController = recommendController{}

func NewRecommendController() *recommendController {
	return &onceRecommendController
}

func (rc *recommendController) TestRecommendList(ctx *fasthttp.RequestCtx) {
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	req := &model.FeedRequest{
		Ctx:      reqCtx,
		Username: "554285007@qq.com",
	}
	err := feed.FeedList(req)
	if err != nil {
		WriterResp(ctx, NewRetData(500, "error", err).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(0, "success", req.FeedResults).ToJson())
}

func (rc *recommendController) RecommendList(ctx *fasthttp.RequestCtx) {
	user := getSessionUser(ctx)
	anyUser, newUser := false, true
	if user == nil {
		anyUser = true
	}
	sessionID := ""
	username := ""
	if anyUser {
		if val := GetSessionID(ctx); val != "" {
			newUser = false
			sessionID = val
		} else {
			sessionID = common.UUID()
		}
		username = fmt.Sprintf("%s-%s", common.AnyUser, sessionID)
	} else {
		username = user.Username
		newUser = false
	}
	reqCtx := RequestContextAndUser(ctx, username)
	defer DestroyContext(reqCtx)
	request := &model.FeedRequest{
		AnyUser:  anyUser,
		Username: username,
		Ctx:      reqCtx,
	}
	err := feed.FeedList(request)
	if err != nil {
		WriterResp(ctx, NewRetData(500, "error", err).ToJson())
		return
	}
	ret := dto.RecommendResults{
		Results: request.FeedResults,
	}
	if newUser {
		ret.SessionID = sessionID
	}
	WriterResp(ctx, NewRetData(successCode, success, ret).ToJson())
}
