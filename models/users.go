package models

import (
	"errors"
	"time"

	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

const PasswordCost = 12

type User struct {
	ID        int       `db:"id,omitempty"`
	CreatedAT time.Time `db:"created_at"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Password  string    `db:"password_hash"`
	Activated bool      `db:"activated"`
}

type UsersModel struct {
	db db.Session
}

// return collection name
func (u UsersModel) Table() string {
	return "users"
}

// gets user from database by id
func (u UsersModel) Get(id int) (*User, error) {
	var user User

	collection := u.db.Collection(u.Table())
	err := collection.Find(db.Cond{"id": id}).One(&user)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrNoMoreRows
		}
	}

	return &user, nil
}

// finds user by email from database
func (u UsersModel) FindByEmail(email string) (*User, error) {
	var user User

	collection := u.db.Collection(u.Table())
	err := collection.Find(db.Cond{"email": email}).One(&user)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrNoMoreRows
		}
	}

	return &user, nil
}

// inserts a user into database
func (u UsersModel) Insert(user *User) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), PasswordCost)
	if err != nil {
		return err
	}

	user.Password = string(newHash)
	user.CreatedAT = time.Now()

	collection := u.db.Collection(u.Table())
	res, err := collection.Insert(user)
	if err != nil {
		switch {
		case errHasDuplicates(err, "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	user.ID = convertUpperIDtoInt(res.ID())

	return nil
}

// Authenticates by checking whether email exists in database and cop
func (u UsersModel) Authenticate(email, password string) (*User, error) {
	user, err := u.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.Activated {
		return nil, ErrUserNotActive
	}

	match, err := user.ComparePassword(password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrInvalidLogin
	}

	return user, nil
}

// compares hashed user password with plainPassword
func (u User) ComparePassword(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
