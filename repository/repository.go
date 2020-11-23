package repository

import "database/sql"

var database *sql.DB

var (
	insertVideoStatement                            *sql.Stmt
	selectVideosWithCaptionAndNotProcessedStatement *sql.Stmt
	updateVideoStatement                            *sql.Stmt
	selectVideosWhereCaptionLike                    *sql.Stmt

	insertCaptionStatement     *sql.Stmt
	selectCaptionsWhereVideoId *sql.Stmt

	insertFrameStatement          *sql.Stmt
	selectFramesByVideoWithinTime *sql.Stmt
	selectFrameById               *sql.Stmt
)

func SetDatabase(db *sql.DB) {
	database = db

	insertVideoStatement, _ = database.Prepare("INSERT INTO videos (id, video_id, caption, has_caption, processed) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;")
	selectVideosWithCaptionAndNotProcessedStatement, _ = database.Prepare("SELECT id, video_id, caption, has_caption, processed FROM videos WHERE has_caption = true AND processed = false;")
	updateVideoStatement, _ = database.Prepare("UPDATE videos set caption = $2, processed = $3  WHERE id = $1;")
	selectVideosWhereCaptionLike, _ = database.Prepare("SELECT id, video_id, caption, has_caption, processed  FROM videos WHERE lower(caption) LIKE $1;")

	insertCaptionStatement, _ = database.Prepare(`INSERT INTO captions (id, video_id, text, start, "end", previous_caption, next_caption) VALUES ($1, $2, $3, $4, $5, $6, $7);`)
	selectCaptionsWhereVideoId, _ = database.Prepare(`SELECT id, text, start, "end", previous_caption, next_caption FROM captions WHERE video_id = $1 ORDER BY start`)

	insertFrameStatement, _ = database.Prepare("INSERT INTO frames (id, video_id, time, image, previous_frame, file_name) VALUES ($1, $2, $3, $4, $5, $6);")
	selectFramesByVideoWithinTime, _ = database.Prepare("SELECT id, time FROM frames WHERE video_id = $1 AND time >= $2 AND time <= $3;")
	selectFrameById, _ = database.Prepare("SELECT image FROM frames WHERE id = $1;")
}

func Begin() (*sql.Tx, error) {
	return database.Begin()
}
