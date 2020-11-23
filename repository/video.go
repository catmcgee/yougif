package repository

import (
	"catmcgee/model"
	"context"
	"database/sql"
	"strings"
)

func InsertVideo(ctx context.Context, video *model.Video) error {
	if _, err := insertVideoStatement.ExecContext(ctx, video.Id, video.VideoId, video.Caption, video.HasCaption, video.Processed); err != nil {
		return err
	}
	return nil
}

func SelectCaptionAndNotProcessedVideos(ctx context.Context) ([]*model.Video, error) {
	rows, err := selectVideosWithCaptionAndNotProcessedStatement.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToVideos(rows)
}

func UpdateVideo(ctx context.Context, tx *sql.Tx, video *model.Video) error {
	stmt := tx.Stmt(updateVideoStatement)
	if _, err := stmt.ExecContext(ctx, video.Id, video.Caption, video.Processed); err != nil {
		return err
	}
	return nil
}

func SelectVideosWhereCaptionLike(ctx context.Context, searchString string) ([]*model.Video, error) {

	rows, err := selectVideosWhereCaptionLike.QueryContext(ctx, "%"+strings.ToLower(searchString)+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToVideos(rows)
}

func rowsToVideos(rows *sql.Rows) ([]*model.Video, error) {
	var id, videoId, caption string
	var hasCaption, processed bool

	result := make([]*model.Video, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &videoId, &caption, &hasCaption, &processed); err != nil {
			return nil, err
		}

		result = append(result, &model.Video{
			Id:         id,
			VideoId:    videoId,
			Caption:    caption,
			HasCaption: hasCaption,
			Processed:  processed,
		})
	}
	return result, nil
}
