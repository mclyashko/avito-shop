package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		JwtSecretKey: []byte("test_secret"),
	}

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{
			name:         "No Authorization header",
			authHeader:   "",
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			name:         "Invalid format",
			authHeader:   "Basic token",
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalidtoken",
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			name:         "Expired token",
			authHeader:   "Bearer " + createJWT(cfg.JwtSecretKey, time.Now().Add(-time.Hour)),
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			name:         "Valid token",
			authHeader:   "Bearer " + createJWT(cfg.JwtSecretKey, time.Now().Add(time.Hour)),
			expectedCode: fiber.StatusOK,
		},
	}

	app := fiber.New()
	app.Use(AuthMiddleware(cfg))
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }) // Фейковый маршрут

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest("GET", "/", tt.authHeader, nil)
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func newTestRequest(method, path, authHeader string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	return req
}

func createJWT(secretKey []byte, exp time.Time) string {
	claims := &service.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)
	return signedToken
}
