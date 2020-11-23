package model

import "github.com/google/uuid"

type Video struct {
	Id         string
	VideoId    string
	Caption    string
	HasCaption bool
	Processed  bool
}

func NewVideo(videoId, caption string, hasCaption bool) *Video {
	return &Video{
		Id:         uuid.New().String(),
		VideoId:    videoId,
		Caption:    caption,
		HasCaption: hasCaption,
	}
}
