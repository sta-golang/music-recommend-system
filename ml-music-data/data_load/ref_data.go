package data_load

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
	PublishTime int    `json:"publishTime"`
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
