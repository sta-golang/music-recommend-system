package controller

import (
	er "github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/model"
	"github.com/sta-golang/music-recommend/service"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type creatorController struct {
}

var onceCreatorController = creatorController{}

func NewCreatorController() *creatorController {
	return &onceCreatorController
}

// GetCreator 获取作者详情
func (cc *creatorController) GetCreator(ctx *fasthttp.RequestCtx) {
	idStr := ctx.FormValue("id")
	if idStr == nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	id, err := strconv.Atoi(str.BytesToString(idStr))
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, paramsErrMessage).ToJson())
		return
	}
	detail, sErr := service.PubCreatorService.GetCreatorDetail(id)
	if sErr != nil && sErr.Err != nil {
		WriterResp(ctx, NewRetDataForErrAndMessage(sErr, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, detail).ToJson())
}

// GetCreators 获取作者列表
func (cc *creatorController) GetCreators(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	page := args.GetUintOrZero("page")
	tp := args.GetUintOrZero("type")
	var creators []model.Creator
	var sErr *er.Error
	if tp == 0 {
		creators, sErr = service.PubCreatorService.GetCreator(page)
		if sErr != nil && sErr.Err != nil {
			WriterResp(ctx, NewRetDataForErrAndMessage(sErr, serverSelectErrMessage).ToJson())
			return
		}
	} else {
		creators, sErr = service.PubCreatorService.GetCreatorWithType(tp, page)
		if sErr != nil && sErr.Err != nil {
			WriterResp(ctx, NewRetDataForErrAndMessage(sErr, serverSelectErrMessage).ToJson())
			return
		}
	}
	WriterResp(ctx, NewRetData(successCode, success, creators).ToJson())
}
