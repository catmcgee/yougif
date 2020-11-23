package model

import "github.com/google/uuid"

type Caption struct {
	Id              string
	VideoId         string
	Text            string
	Start           int
	End             int
	PreviousCaption *string
	NextCaption     *string
}

func NewCaption(videoId, text string, start, end int) *Caption {
	return &Caption{
		Id:              uuid.New().String(),
		VideoId:         videoId,
		Text:            text,
		Start:           start,
		End:             end,
		PreviousCaption: nil,
		NextCaption:     nil,
	}
}
