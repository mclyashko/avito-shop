package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/db"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var ErrWrongPassword = fmt.Errorf("wrong password")

func GetTokenByUsernameAndPassword(ctx context.Context, pool *pgxpool.Pool, username string, password string, jwtSecretKey []byte) (*string, error) {
	user, err := db.GetUserByLogin(ctx, pool, username)
	if err == db.ErrUserNotFound {
		user, err = db.CreateNewUser(ctx, pool, username, password)
	}
	if err != nil {
		return nil, err
	}

	if !comparePasswords(user.PasswordHash, password) {
		log.Printf("Passwords are not equal: %v %v", user.PasswordHash, password)
		return nil, ErrWrongPassword
	}

	claims := &JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecretKey)

	if err != nil {
		log.Printf("Error signing token: %v", err)
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	return &signedToken, nil
}

func comparePasswords(hashedPassword, password string) bool {
	hash := sha256.New()
	hash.Write([]byte(password))
	enteredPasswordHash := hex.EncodeToString(hash.Sum(nil))

	return hashedPassword == enteredPasswordHash
}
