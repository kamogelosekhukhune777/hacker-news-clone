package models

import "github.com/upper/db/v4"

type Models struct {
	Users UsersModel
}

func New(db db.Session) Models {
	return Models{
		Users: UsersModel{db: db},
	}
}
