package data_load

type WangYiYunResult struct {
	Result RefResult `json:"result"`
}

type RefResult struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Name string `json:"name"`
	ID int `json:"id"`
	Artists []Artist `json:"artists"`
	Album Album `json:"album"`
	Duration int `json:"duration"`
}

type Artist struct {
	Name string `json:"name"`
	ID int `json:"id"`

}
