package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/model"
)

type TestUserAccessor struct {
	users map[string]*model.User
	db.UserAccessor
}

func (t *TestUserAccessor) GetUserByLogin(ctx context.Context, username string) (*model.User, error) {
	user, exists := t.users[username]
	if !exists {
		return nil, db.ErrUserNotFound
	}
	return user, nil
}

func (t *TestUserAccessor) InsertNewUser(ctx context.Context, username string, passwordHash string, balance int64) (*model.User, error) {
	user := &model.User{
		Login:        username,
		PasswordHash: passwordHash,
		Balance:      balance,
	}
	t.users[username] = user
	return user, nil
}

func TestAuthServiceImpl_GetTokenByUsernameAndPassword(t *testing.T) {
	config := &config.Config{
		JwtSecretKey:          []byte("12345"),
		JwtExpirationDuration: time.Hour,
	}
	s := &basicServiceImpl{
		cfg: config,
	}
	testUserAccessor := &TestUserAccessor{users: make(map[string]*model.User)}
	authService := &AuthServiceImpl{
		Service:      s,
		UserAccessor: testUserAccessor,
	}

	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantToken bool
	}{
		{
			name: "create user",
			args: args{
				ctx:      context.Background(),
				username: "user",
				password: "password",
			},
			wantErr:   false,
			wantToken: true,
		},
		{
			name: "get user",
			args: args{
				ctx:      context.Background(),
				username: "user",
				password: "password",
			},
			wantErr:   false,
			wantToken: true,
		},
		{
			name: "get user with wrong password",
			args: args{
				ctx:      context.Background(),
				username: "user",
				password: "password123",
			},
			wantErr:   true,
			wantToken: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authService.GetTokenByUsernameAndPassword(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServiceImpl.GetTokenByUsernameAndPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantToken {
				if got == nil {
					t.Errorf("AuthServiceImpl.GetTokenByUsernameAndPassword() = nil, want token")
				} else {
					tokenString := *got
					parsedToken, parseErr := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
						return config.JwtSecretKey, nil
					})
					if parseErr != nil {
						t.Errorf("Error parsing token: %v", parseErr)
						return
					}

					if claims, ok := parsedToken.Claims.(*JWTClaims); ok && parsedToken.Valid {
						if claims.Username != tt.args.username {
							t.Errorf("Expected username in token to be %v, but got %v", tt.args.username, claims.Username)
						}
					} else {
						t.Errorf("Invalid token or invalid claims")
					}
				}
			}
		})
	}
}

func Test_comparePasswords(t *testing.T) {
	type args struct {
		hashedPassword string
		password       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Correct match for password123",
			args: args{
				hashedPassword: "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
				password:       "password123",
			},
			want: true,
		},
		{
			name: "Correct match for qwerty",
			args: args{
				hashedPassword: "65e84be33532fb784c48129675f9eff3a682b27168c0ea744b2cf58ee02337c5",
				password:       "qwerty",
			},
			want: true,
		},
		{
			name: "Hashes are not equal",
			args: args{
				hashedPassword: "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94e",
				password:       "password123",
			},
			want: false,
		},
		{
			name: "Password is empty",
			args: args{
				hashedPassword: "ef92b778bafee8f392d76a89a5c70ff8183f46b11fcde79c521ddf0a4b35dd54",
				password:       "",
			},
			want: false,
		},
		{
			name: "Hash from empty password",
			args: args{
				hashedPassword: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				password:       "password123",
			},
			want: false,
		},
		{
			name: "Hash from empty password and password is empty",
			args: args{
				hashedPassword: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				password:       "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := comparePasswords(tt.args.hashedPassword, tt.args.password); got != tt.want {
				t.Errorf("comparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}
