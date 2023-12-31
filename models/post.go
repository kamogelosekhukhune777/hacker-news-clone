package models

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/upper/db/v4"
)

var (
	ErrDuplicateTitle = errors.New("title already exist in database")
	ErrDuplicateVotes = errors.New("you already voted")

	queryTemplate = `SELECT COUNT(*) OVER() AS total_records, pq.*, u.name as uname FROM (
	    SELECT p.id, p.title, p.url, p.created_at, p.user_id as uid, COUNT(c.post_id) as comment_count, count(v.post_id) as votes
		FROM posts p 
		LEFT JOIN comments c ON p.id = c.post_id 
	    LEFT JOIN votes v ON p.id = v.post_id
	 	#where#
		GROUP BY p.id
		#orderby#
		) AS pq
	LEFT JOIN users u ON u.id = uid
	#limit#`
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

func (p PostModel) Table() string {
	return "posts"
}

// gets a specific post
func (p PostModel) Get(id int) (*Post, error) {
	var post Post
	q := strings.Replace(queryTemplate, "#where#", "WHERE p.id = $1", 1)
	q = strings.Replace(q, "#orderby#", "", 1)
	q = strings.Replace(q, "#limit#", "", 1)

	row, err := p.db.SQL().Query(q, id)
	if err != nil {
		return nil, err
	}

	iterator := p.db.SQL().NewIterator(row)
	err = iterator.One(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// gets all epecific posts
func (p PostModel) GetAll(f Filter) ([]Post, Metadata, error) {
	var posts []Post
	var rows *sql.Rows
	var err error
	meta := Metadata{}

	q := f.applyTemplate(queryTemplate)
	if len(f.Query) > 0 {
		rows, err = p.db.SQL().Query(q, "%"+strings.ToLower(f.Query)+"%", f.limit(), f.offset())
	} else {
		rows, err = p.db.SQL().Query(q, f.limit(), f.offset())
	}

	if err != nil {
		return nil, meta, err
	}

	iter := p.db.SQL().NewIterator(rows)
	err = iter.All(&posts)

	if err != nil {
		return nil, meta, err
	}

	if len(posts) == 0 {
		return nil, meta, nil
	}

	first := posts[0]
	return posts, calculateMetadata(first.TotalRecords, f.Page, f.PageSize), nil
}

func (p PostModel) Vote(postId, userId int) error {

	collection := p.db.Collection("votes")

	_, err := collection.Insert(map[string]int{
		"post_id": postId,
		"user_id": userId,
	})
	if err != nil {
		if errHasDuplicates(err, "votes_pkey") {
			return ErrDuplicateVotes
		}
		return err
	}

	return nil
}

func (p *Post) DateHuman() string {
	return carbon.CreateFromStdTime(p.CreatedAt).DiffForHumans()
}

func (p *Post) Host() string {
	ur, err := url.Parse(p.Url)
	if err != nil {
		return ""
	}

	return ur.Host
}

func (p PostModel) Insert(title, url string, userId int) (*Post, error) {

	post := Post{
		CreatedAt: time.Now(),
		Title:     title,
		Url:       url,
		UserId:    userId,
	}

	collection := p.db.Collection(p.Table())

	res, err := collection.Insert(post)
	if err != nil {
		return nil, err
	}

	post.ID = convertUpperIDtoInt(res.ID())

	return &post, nil
}
