package service

import (
	"catmcgee/model"
	"catmcgee/repository"
	"context"
	"github.com/sirupsen/logrus"
)

func SearchForVideo(ctx context.Context, searchString string) (*model.SearchResult, error) {
	videos, err := repository.SelectVideosWhereCaptionLike(ctx, searchString)
	if err != nil {
		return nil, err
	}

	result := make([]*model.FrameResult, 0)
	for _, video := range videos {
		captions, err := repository.SelectCaptionsWhereVideoId(ctx, video.Id)
		if err != nil {
			logrus.Println(err)
			continue
		}

		captionsWithSearchString := getCaptionsWithSearchString(captions, searchString)

		for _, caption := range captionsWithSearchString {
			frames, err := repository.SelectByVideoAndWithinTimeFrames(ctx, video.Id, caption.Start, caption.End)
			if err != nil {
				logrus.Println(err)
				continue
			}

			for _, frame := range frames {
				result = append(result, &model.FrameResult{
					Id:             frame.Id,
					YouTubeVideoId: video.VideoId,
					Time:           frame.Time,
				})
			}

		}
	}
	return &model.SearchResult{Data: result}, nil
}
