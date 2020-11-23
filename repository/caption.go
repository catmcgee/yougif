package repository

import (
	"catmcgee/model"
	"context"
	"database/sql"
)

func InsertCaption(ctx context.Context, tx *sql.Tx, caption *model.Caption) error {
	stmt := tx.Stmt(insertCaptionStatement)
	if _, err := stmt.ExecContext(ctx, caption.Id, caption.VideoId, caption.Text, caption.Start, caption.End, caption.PreviousCaption, caption.NextCaption); err != nil {
		return err
	}
	return nil
}

func SelectCaptionsWhereVideoId(ctx context.Context, videoId string) ([]*model.Caption, error) {
	rows, err := selectCaptionsWhereVideoId.QueryContext(ctx, videoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Caption, 0)

	var id, text string
	var previousCaption, nextCaption *string
	var start, end int

	for rows.Next() {
		if err := rows.Scan(&id, &text, &start, &end, &previousCaption, &nextCaption); err != nil {
			return nil, err
		}

		result = append(result, &model.Caption{
			Id:              id,
			VideoId:         videoId,
			Text:            text,
			Start:           start,
			End:             end,
			PreviousCaption: previousCaption,
			NextCaption:     nextCaption,
		})
	}

	return result, nil
}
