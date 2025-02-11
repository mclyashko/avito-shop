package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mclyashko/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct{}

func (s *MockAuthService) GetTokenByUsernameAndPassword(ctx context.Context, username, password string) (*string, error) {
	if username == "testuser" && password == "testpass" {
		token := "validToken123"
		return &token, nil
	}
	if username == "testuser" && password == "wrongpass" {
		return nil, service.ErrWrongPassword
	}
	return nil, service.ErrWrongPassword
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name         string
		body         []byte
		mockResponse string
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful authentication",
			body:         []byte(`{"username":"testuser", "password":"testpass"}`),
			mockResponse: "validToken123",
			mockError:    nil,
			expectedCode: fiber.StatusOK,
			expectedBody: `{"token":"validToken123"}`,
		},
		{
			name:         "Wrong password",
			body:         []byte(`{"username":"testuser", "password":"wrongpass"}`),
			mockResponse: "",
			mockError:    service.ErrWrongPassword,
			expectedCode: fiber.StatusUnauthorized,
			expectedBody: `{"errors":"Wrong password"}`,
		},
		{
			name:         "Invalid request body",
			body:         []byte(`invalid json`),
			mockResponse: "",
			mockError:    nil,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Invalid request"}`,
		},
	}

	app := fiber.New()

	mockAuthService := &MockAuthService{}

	app.Post("/api/auth", func(c *fiber.Ctx) error {
		return Authenticate(c, mockAuthService)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest("POST", "/api/auth", tt.body)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error while making the request: %v", err)
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var responseBody map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&responseBody)
			responseJSON, _ := json.Marshal(responseBody)

			assert.JSONEq(t, tt.expectedBody, string(responseJSON))
		})
	}
}

func newTestRequest(method, path string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
