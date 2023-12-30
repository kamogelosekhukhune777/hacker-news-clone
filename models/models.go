package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/upper/db/v4"
)

var (
	ErrNoMoreRows     = errors.New("no record found")
	ErrDuplicateEmail = errors.New("email already exists in database")
	ErrUserNotActive  = errors.New("your account is inactive")
	ErrInvalidLogin   = errors.New("Invalid login")
)

type Models struct {
	Users UsersModel
}

func New(db db.Session) Models {
	return Models{
		Users: UsersModel{db: db},
	}
}

func convertUpperIDtoInt(id db.ID) int {
	idType := fmt.Sprintf("%T", id)
	if idType == "int64" {
		return int(id.(int64))
	}

	return id.(int)
}

func errHasDuplicates(err error, key string) bool {
	str := fmt.Sprintf(`ERROR: duplicate key value violates unique constraint "%s"`, key)
	return strings.Contains(err.Error(), str)
}
