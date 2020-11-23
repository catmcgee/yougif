package model

type Captions struct {
	CaptionTracks []CaptionTrack `json:"captionTracks"`
}

type CaptionTrack struct {
	BaseUrl      string `json:"baseUrl"`
	VssId        string `json:"vssId"`
	LanguageCode string `json:"languageCode"`
}
