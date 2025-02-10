package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	log.Printf("Got user: %v", user)
	return &user, nil
}

func CreateNewUser(ctx context.Context, pool *pgxpool.Pool, login string, password string) (*model.User, error) {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	user := model.User{
		Login: login,
		PasswordHash: hashedPassword,
		Balance: initialBalance,
	}

	query := `INSERT INTO "user" (login, password_hash, balance) VALUES ($1, $2, $3)`

	_, err := pool.Exec(ctx, query, user.Login, user.PasswordHash, user.Balance)

	if err != nil {
		log.Printf("Failed to create new user: %v", err)
		return nil, fmt.Errorf("could not create user: %w", err)
	}


	log.Printf("User created: %v", user)
	return &user, nil
}