package model

import "github.com/google/uuid"

type Frame struct {
	Id            string
	VideoId       string
	FileName      string
	Time          int
	Image         []byte
	PreviousFrame string
}

func NewFrame(videoId, fileName string, time int, image []byte, previousFrame string) *Frame {
	return &Frame{
		Id:            uuid.New().String(),
		VideoId:       videoId,
		FileName:      fileName,
		Time:          time,
		Image:         image,
		PreviousFrame: previousFrame,
	}
}
