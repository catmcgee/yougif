package repository

import (
	"catmcgee/model"
	"context"
	"database/sql"
)

func InsertFrame(ctx context.Context, tx *sql.Tx, frame *model.Frame) error {
	stmt := tx.Stmt(insertFrameStatement)
	if _, err := stmt.ExecContext(ctx, frame.Id, frame.VideoId, frame.Time, frame.Image, frame.PreviousFrame, frame.FileName); err != nil {
		return err
	}
	return nil
}

func SelectByVideoAndWithinTimeFrames(ctx context.Context, videoId string, start, end int) ([]*model.Frame, error) {
	rows, err := selectFramesByVideoWithinTime.QueryContext(ctx, videoId, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Frame, 0)

	var id string
	var time int

	for rows.Next() {
		if err := rows.Scan(&id, &time); err != nil {
			return nil, err
		}

		result = append(result, &model.Frame{
			Id:   id,
			Time: time,
		})
	}

	return result, nil
}

func SelectByIdFrame(ctx context.Context, id string) ([]byte, error) {
	row := selectFrameById.QueryRowContext(ctx, id)

	var image []byte
	if err := row.Scan(&image); err != nil {
		return nil, err
	}

	return image, nil
}
