package controller

import (
	"github.com/sta-golang/go-lib-utils/str"
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
