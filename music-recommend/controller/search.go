package controller

import (
	"net/http"

	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-recommend/service"
	"github.com/valyala/fasthttp"
)

type SearchController struct{}

var onceSearchController = SearchController{}

func NewSearchController() *SearchController {
	return &onceSearchController
}

func (sc *SearchController) SearchKeyWorld(ctx *fasthttp.RequestCtx) {
	keyword := str.BytesToString(ctx.FormValue("keywords"))
	if keyword == "" {
		WriterResp(ctx, NewRetData(successCode, success, nil).ToJson())
		return
	}
	reqCtx := RequestContext(ctx)
	defer DestroyContext(reqCtx)
	ret, err := service.PubSearchService.SearchKeyworld(reqCtx, keyword)
	if err != nil {
		WriterResp(ctx, NewRetDataForErrorAndMessage(http.StatusBadRequest, serverSelectErrMessage).ToJson())
		return
	}
	WriterResp(ctx, NewRetData(successCode, success, ret).ToJson())
}

func (sc *SearchController) SearchHot(ctx *fasthttp.RequestCtx) {
	retData := `{"code":200,"data":[{"searchWord":"My Cookie Can","score":4073392,"content":"听了这首歌 少女心upup","source":0,"iconType":1,"iconUrl":"https://p1.music.126.net/2zQ0d1ThZCX5Jtkvks9aOQ==/109951163968000522.png","url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"乌鸦","score":1844069,"content":"不想合群 但我要和你成群","source":0,"iconType":2,"iconUrl":"https://p1.music.126.net/szWeddITZIVxpvQ0QywzcQ==/109951163967989323.png","url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"还是会想你","score":1842941,"content":"如果，我们变成彼此爱情的局外人","source":0,"iconType":0,"iconUrl":"http://p4.music.126.net/P4mXkx6VKXLFqVo5ohHxDg==/109951163992439900.jpg","url":"","alg":"featured"},{"searchWord":"lovely","score":1835362,"content":"歌名叫lovely 感觉到lonely","source":0,"iconType":5,"iconUrl":"https://p1.music.126.net/PtgUJbcvDAY0HKWpmOB2HA==/109951163967988470.png","url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"海底","score":1546686,"content":"温柔的人会将你带离海底","source":0,"iconType":1,"iconUrl":"https://p1.music.126.net/2zQ0d1ThZCX5Jtkvks9aOQ==/109951163968000522.png","url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"失重","score":1352268,"content":"听说这首歌能缓解65%的焦虑","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"日不落","score":1107831,"content":"我要送你 日不落的爱恋","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"许嵩","score":1062941,"content":"许嵩的歌，雅俗都能共赏~","source":0,"iconType":1,"iconUrl":"https://p1.music.126.net/2zQ0d1ThZCX5Jtkvks9aOQ==/109951163968000522.png","url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"红马","score":975866,"content":"你用斜阳揉碎了春雪","source":0,"iconType":0,"iconUrl":"http://p3.music.126.net/P4mXkx6VKXLFqVo5ohHxDg==/109951163992439900.jpg","url":"","alg":"featured"},{"searchWord":"薛之谦","score":931326,"content":"老薛一发歌就能掀起狂潮！","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"起风了","score":884040,"content":"心之所动 且就随缘去吧","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"沦陷","score":793995,"content":"影视剧《原来我很爱你》片尾曲","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"Mood","score":760416,"content":"I'm not in a mood","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"张哲瀚","score":753253,"content":"你知道，他本身就是光。","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"张杰","score":752679,"content":"张杰带你穿越人海","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"童话镇","score":707753,"content":"听说白雪公主在逃跑","source":0,"iconType":5,"iconUrl":"https://p1.music.126.net/PtgUJbcvDAY0HKWpmOB2HA==/109951163967988470.png","url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"溯","score":707708,"content":"慵懒歌曲，失眠的首选","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"林俊杰","score":670472,"content":"当之无愧的行走CD机！","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"错位时空","score":617697,"content":"我吹过你吹过的晚风，那我们算不算相拥","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"},{"searchWord":"毛不易","score":612657,"content":"深情的唱作人毛不易！","source":0,"iconType":0,"iconUrl":null,"url":"","alg":"alg_search_rec_hotquery_base_hotquery"}],"message":"success"}`
	WriterResp(ctx, str.StringToBytes(&retData))
}
