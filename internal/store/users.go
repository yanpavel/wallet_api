package store

import (
	"context"
	"database/sql"
	"errors"
	"sync"
)

type UsersStore struct {
	db *sql.DB
	mx sync.RWMutex
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UsersStore) GetUser(ctx context.Context, login string) (*User, error) {
	u.mx.RLock()
	defer u.mx.RUnlock()

	query := `
		SELECT id, username, password FROM users WHERE username=$1;
	`

	var user User
	err := u.db.QueryRowContext(ctx, query, login).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UsersStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	u.mx.RLock()
	defer u.mx.RUnlock()

	query := `
		SELECT id, username, password FROM users WHERE id=$1
	`

	var user User
	err := u.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UsersStore) CreateUser(ctx context.Context, login string, password string) (*int64, error) {
	u.mx.Lock()
	defer u.mx.Unlock()

	query := `
	INSERT INTO users (username, password)
	VALUES ($1, $2) RETURNING id;
	`
	var user User
	err := u.db.QueryRowContext(ctx, query, login, password).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return &user.Id, nil
}
