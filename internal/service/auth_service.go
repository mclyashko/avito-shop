package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mclyashko/avito-shop/internal/db"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var ErrWrongPassword = fmt.Errorf("wrong password")

type AuthService interface {
	GetTokenByUsernameAndPassword(ctx context.Context, username string, password string) (*string, error)
}

type AuthServiceImpl struct {
	*Service
	UserAccessor db.UserAccessor
}

const (
	initialBalance = 1000
)

func (s *AuthServiceImpl) GetTokenByUsernameAndPassword(ctx context.Context, username string, password string) (*string, error) {
	user, err := s.UserAccessor.GetUserByLogin(ctx, username)
	if err == db.ErrUserNotFound {
		hash := sha256.New()
		hash.Write([]byte(password))
		hashedPassword := hex.EncodeToString(hash.Sum(nil))

		user, err = s.UserAccessor.InsertNewUser(ctx, username, hashedPassword, initialBalance)
	}
	if err != nil {
		return nil, err
	}

	if !comparePasswords(user.PasswordHash, password) {
		return nil, ErrWrongPassword
	}

	claims := &JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JwtExpirationDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.cfg.JwtSecretKey)

	if err != nil {
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
