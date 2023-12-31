package models

import (
	"time"

	"github.com/upper/db/v4"
)

type Post struct {
	ID           int       `db:"id,omitempty"`
	Title        string    `db:"title"`
	Url          string    `db:"url"`
	CreatedAt    time.Time `db:"created_at"`
	UserId       int       `db:"user_id"`
	Votes        int       `db:"votes"`
	UserName     string    `db:"user_name,omitempty"`
	CommentCount int       `db:"comment_count,omitempty"`
	TotalRecords int       `db:"total_records,omitempty"`
}

type PostModel struct {
	db db.Session
}

//01:14:00
