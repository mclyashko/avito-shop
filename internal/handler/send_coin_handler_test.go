package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

type MockSendCoinsService struct {
	mockError error
}

func (s *MockSendCoinsService) SendCoins(ctx context.Context, fromUser, toUser string, amount int64) error {
	return s.mockError
}

func TestSendCoinHandler(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		body         string
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful transaction",
			token:        "testUser",
			body:         `{"toUser": "receiverUser", "amount": 50}`,
			mockError:    nil,
			expectedCode: fiber.StatusOK,
			expectedBody: ``,
		},
		{
			name:         "No token",
			token:        "",
			body:         `{"toUser": "receiverUser", "amount": 50}`,
			mockError:    nil,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"No token"}`,
		},
		{
			name:         "Invalid JSON",
			token:        "testUser",
			body:         `{"toUser": "receiverUser", "amount": "invalid"}`,
			mockError:    nil,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Invalid request"}`,
		},
		{
			name:         "Negative transaction",
			token:        "testUser",
			body:         `{"toUser": "receiverUser", "amount": -10}`,
			mockError:    service.ErrNegativeSignTransaction,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Negative sign transaction"}`,
		},
		{
			name:         "Receiver not found",
			token:        "testUser",
			body:         `{"toUser": "unknownUser", "amount": 10}`,
			mockError:    db.ErrUserNotFound,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Reciever not found"}`,
		},
		{
			name:         "Insufficient funds",
			token:        "testUser",
			body:         `{"toUser": "receiverUser", "amount": 1000}`,
			mockError:    service.ErrInsufficientFunds,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Insufficient funds"}`,
		},
		{
			name:         "Internal server error",
			token:        "testUser",
			body:         `{"toUser": "receiverUser", "amount": 50}`,
			mockError:    errors.New("database failure"),
			expectedCode: fiber.StatusInternalServerError,
			expectedBody: `{"errors":"Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(AddClaimsToContextFromRequest())

			mockService := &MockSendCoinsService{mockError: tt.mockError}

			app.Post("/api/send", func(c *fiber.Ctx) error {
				return SendCoinHandler(c, mockService)
			})

			req, _ := http.NewRequest("POST", "/api/send", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedBody != "" {
				var responseBody map[string]interface{}
				_ = json.NewDecoder(resp.Body).Decode(&responseBody)
				responseJSON, _ := json.Marshal(responseBody)
				assert.JSONEq(t, tt.expectedBody, string(responseJSON))
			}
		})
	}
}
