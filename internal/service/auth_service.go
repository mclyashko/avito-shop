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
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var ErrWrongPassword = fmt.Errorf("wrong password")

func GetTokenByUsernameAndPassword(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, username string, password string, jwtSecretKey []byte) (*string, error) {
	user, err := db.GetUserByLogin(ctx, pool, username)
	if err == db.ErrUserNotFound {
		hash := sha256.New()
		hash.Write([]byte(password))
		hashedPassword := hex.EncodeToString(hash.Sum(nil))

		user, err = db.InsertNewUser(ctx, pool, username, hashedPassword)
	}
	if err != nil {
		return nil, err
	}

	if !comparePasswords(user.PasswordHash, password) {
		log.Printf("Passwords are not equal for user %s", username)
		return nil, ErrWrongPassword
	}

	claims := &JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JwtExpirationDuration)),
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
