package models

import (
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/upper/db/v4"
)

type CommentsModel struct {
	db db.Session
}

type Comment struct {
	ID        int       `db:"comment_id,omitempty"`
	CreatedAt time.Time `db:"c_created_at,omitempty"`
	Body      string    `db:"body"`
	PostID    int       `db:"post_id"`
	UserID    int       `db:"user_id"`
	User      `db:",inline"`
}

func (c CommentsModel) Table() string {
	return "posts"
}

func (c CommentsModel) GetForPost(postId int) ([]Comment, error) {
	var comments []Comment

	q := c.db.SQL().Select("c.id as comment_id", "c.created_at as c_created_at", "*").
		From("comments as c").
		Join("users as u").On("c.user_id = u.id").
		Where(db.Cond{"c.post_id": postId}).
		OrderBy("c.created_at desc")

	err := q.All(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c CommentsModel) Insert(body string, postId, userId int) error {
	_, err := c.db.Collection(c.Table()).Insert(map[string]interface{}{
		"body":    body,
		"user_id": userId,
		"post_id": postId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *Comment) DateHuman() string {
	return carbon.CreateFromStdTime(c.CreatedAt).DiffForHumans()
}
