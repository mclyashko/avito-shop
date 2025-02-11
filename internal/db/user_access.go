package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/model"
)

const (
	initialBalance = 1000
)

var ErrUserNotFound = fmt.Errorf("user not found")

func GetUserByLogin(ctx context.Context, pool *pgxpool.Pool, login string) (*model.User, error) {
	query := `SELECT login, password_hash, balance FROM "user" WHERE login = $1`

	row := pool.QueryRow(ctx, query, login)

	var user model.User

	err := row.Scan(&user.Login, &user.PasswordHash, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		log.Printf("Failed to get user by login: %v", err)
		return nil, err
	}

	log.Printf("Got user with login: %v", user.Login)
	return &user, nil
}

func GetUserByLoginTx(ctx context.Context, tx pgx.Tx, login string) (*model.User, error) {
	query := `SELECT login, password_hash, balance FROM "user" WHERE login = $1`

	row := tx.QueryRow(ctx, query, login)

	var user model.User

	err := row.Scan(&user.Login, &user.PasswordHash, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		log.Printf("Failed to get user by login: %v", err)
		return nil, err
	}

	log.Printf("Got user with login: %v", user.Login)
	return &user, nil
}


func InsertNewUser(ctx context.Context, pool *pgxpool.Pool, login string, password string) (*model.User, error) {
	user := model.User{
		Login: login,
		PasswordHash: password,
		Balance: initialBalance,
	}

	query := `INSERT INTO "user" (login, password_hash, balance) VALUES ($1, $2, $3)`

	_, err := pool.Exec(ctx, query, user.Login, user.PasswordHash, user.Balance)

	if err != nil {
		log.Printf("Failed to create new user with login: %v, error: %v", login, err)
		return nil, fmt.Errorf("could not create user with login: %v", login)
	}

	log.Printf("User with login %v created", login)
	return &user, nil
}

func UpdateUserBalanceTx(ctx context.Context, tx pgx.Tx, login string, amount int64) error {
	query := `
		UPDATE "user" 
		SET balance = balance + $1 WHERE login = $2
	`
	
	_, err := tx.Exec(ctx, query, amount, login)
	if err != nil {
		log.Printf("Failed to update balance for login: %v, error: %v", login, err)
		return err
	}
	
	return nil
}
