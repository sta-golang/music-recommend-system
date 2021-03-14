package data_load

import "fmt"

const (
	APICrawlerName = "http://127.0.0.1:3000"
	WangYiMusic    = "网易音乐人"
)

type WangYiYunResult struct {
	Result RefResult `json:"result"`
}

type RefResult struct {
	Tracks []Track  `json:"tracks"`
	Tags   []string `json:"tags"`
}

type Track struct {
	Name     string   `json:"name"`
	ID       int      `json:"id"`
	Artists  []Artist `json:"artists"`
	Album    Album    `json:"album"`
	Duration int      `json:"duration"`
}

type Artist struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Album struct {
	Name        string `json:"name"`
	BlurPicUrl  string `json:"blurPicUrl"`
	PublishTime int64  `json:"publishTime"`
}

type CrawlerCreator struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	ImageUrl       string   `json:"image_url"`
	Description    string   `json:"description"`
	Superstar      bool     `json:"superstar"`
	SimilarCreator []string `json:"similar_creator"`
	FansNum        int      `json:"fans_num"`
}

type APICreatorResult struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    APICreatorData `json:"data"`
}

func (a *APICreatorResult) GetUrl(id int) string {
	return fmt.Sprintf("/artist/detail?id=%d", id)
}

type APICreatorData struct {
	Artist APICreatorArtist `json:"artist"`
}

type APICreatorArtist struct {
	ID          int      `json:"id"`
	Cover       string   `json:"cover"`
	Name        string   `json:"name"`
	IdentifyTag []string `json:"identifyTag"`
	BriefDesc   string   `json:"briefDesc"`
}

type APISimilarResult struct {
	Artists []APIArtistData `json:"artists"`
}

func (a *APISimilarResult) GetUrl(id int) string {
	return fmt.Sprintf("/simi/artist?id=%d", id)
}

type APIArtistData struct {
	ID int `json:"id"`
}

type APICreatorListResult struct {
	Artists []APIArtistData `json:"artists"`
	More    bool            `json:"more"`
}

type APICreatorMusicResult struct {
	Code  int              `json:"code"`
	More  bool             `json:"more"`
	Songs []APICreatorSong `json:"songs"`
	Total int              `json:"total"`
}

type APICreatorSong struct {
	ID int `json:"id"` // id
}

type APIMusicDetailResult struct {
	Code  int              `json:"code"`
	Songs []APIMusicDetail `json:"songs"`
}

type APIMusicDetail struct {
	Name        string             `json:"name"`
	ID          int                `json:"id"`
	Dt          int                `json:"dt"`
	PublishTime int64              `json:"publishTime"`
	AR          []APIMusicDetailAR `json:"ar"`
	AL          APIMusicDetailAL   `json:"al"`
}

type APIMusicDetailAR struct {
	CreatorID   int    `json:"id"`
	CreatorName string `json:"name"`
}

type APIMusicDetailAL struct {
	TitleName string `json:"name"`
	TitleUrl  string `json:"picUrl"`
}
