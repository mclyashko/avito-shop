package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/model"
)

var ErrUserNotFound = fmt.Errorf("user not found")

type UserAccessor interface {
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByLoginTx(ctx context.Context, tx pgx.Tx, login string) (*model.User, error)
	InsertNewUser(ctx context.Context, login string, password string, balance int64) (*model.User, error)
	UpdateUserBalanceTx(ctx context.Context, tx pgx.Tx, login string, amount int64) error
}

type UserAccessorImpl struct {
	*Db
}

func (db *UserAccessorImpl) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `SELECT login, password_hash, balance FROM "user" WHERE login = $1`

	row := db.pool.QueryRow(ctx, query, login)

	var user model.User

	err := row.Scan(&user.Login, &user.PasswordHash, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db *UserAccessorImpl) GetUserByLoginTx(ctx context.Context, tx pgx.Tx, login string) (*model.User, error) {
	query := `SELECT login, password_hash, balance FROM "user" WHERE login = $1`

	row := tx.QueryRow(ctx, query, login)

	var user model.User

	err := row.Scan(&user.Login, &user.PasswordHash, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db *UserAccessorImpl) InsertNewUser(ctx context.Context, login string, password string, balance int64) (*model.User, error) {
	user := model.User{
		Login:        login,
		PasswordHash: password,
		Balance:      balance,
	}

	query := `INSERT INTO "user" (login, password_hash, balance) VALUES ($1, $2, $3)`

	_, err := db.pool.Exec(ctx, query, user.Login, user.PasswordHash, user.Balance)

	if err != nil {
		return nil, fmt.Errorf("could not create user with login: %v", login)
	}

	return &user, nil
}

func (db *UserAccessorImpl) UpdateUserBalanceTx(ctx context.Context, tx pgx.Tx, login string, amount int64) error {
	query := `
		UPDATE "user" 
		SET balance = balance + $1 WHERE login = $2
	`

	_, err := tx.Exec(ctx, query, amount, login)
	if err != nil {
		return err
	}

	return nil
}
