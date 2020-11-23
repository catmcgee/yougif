package model

type FrameResult struct {
	Id             string `json:"id"`
	YouTubeVideoId string `json:"youTubeVideoId"`
	Time           int    `json:"timeInMs"`
}

type SearchResult struct {
	Data []*FrameResult `json:"data"`
}
