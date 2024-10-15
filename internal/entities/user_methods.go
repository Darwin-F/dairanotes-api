package entities

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type UserMethodsInterface interface {
	Store(ctx context.Context, newUser User) error
	Update(ctx context.Context, userID int64, user User) error
	Destroy(ctx context.Context, userID int64) error
	GetPasswordByUserName(ctx context.Context, username string) (string, error)
}

type UserMethods struct {
	DB *sqlx.DB
}

func NewUserMethods(db *sqlx.DB) *UserMethods {
	return &UserMethods{
		DB: db,
	}
}

func (u *UserMethods) Store(ctx context.Context, newUser User) error {
	_, err := u.DB.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserMethods) Update(ctx context.Context, userID int64, user User) error {
	_, err := u.DB.ExecContext(ctx, "UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", user.Username, user.Email, user.Password, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserMethods) Destroy(ctx context.Context, userID int64) error {
	_, err := u.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserMethods) GetPasswordByUserName(ctx context.Context, username string) (string, error) {
	var password string
	err := u.DB.QueryRowContext(ctx, "SELECT password FROM users WHERE username = ?", username).Scan(&password)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUserNotFound
	}

	return password, nil
}
