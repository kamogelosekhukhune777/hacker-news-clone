package models

import "github.com/upper/db/v4"

type UsersModel struct {
	db db.Session
}

func (u UsersModel) Get() {
	u.db.Collection("users").Find(db.Cond{"id": 1})
}
