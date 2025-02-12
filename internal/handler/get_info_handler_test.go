package handler

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/model"
	"github.com/mclyashko/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

func AddClaimsToContextFromRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token != "" {
			c.Locals("claims", &service.JWTClaims{
				Username: token,
			})
		}
		return c.Next()
	}
}

type MockUserInfoService struct{}

func (s *MockUserInfoService) GetUserInfo(ctx context.Context, username string) (*int64, []model.UserItem, []model.CoinTransfer, []model.CoinTransfer, error) {
	if username != "testUser" {
		return nil, nil, nil, nil, db.ErrUserNotFound
	}

	coins := int64(100)
	userItems := []model.UserItem{
		{ItemName: "Sword", Quantity: 1},
		{ItemName: "Shield", Quantity: 2},
	}
	receivedTransfers := []model.CoinTransfer{
		{SenderLogin: "Alice", Amount: 50},
	}
	sentTransfers := []model.CoinTransfer{
		{ReceiverLogin: "Bob", Amount: 20},
	}

	return &coins, userItems, receivedTransfers, sentTransfers, nil
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful retrieval",
			token:        "testUser",
			expectedCode: fiber.StatusOK,
			expectedBody: `{
				"coins": 100,
				"inventory": [
					{"type": "Sword", "quantity": 1},
					{"type": "Shield", "quantity": 2}
				],
				"coinHistory": {
					"received": [{"fromUser": "Alice", "amount": 50}],
					"sent": [{"toUser": "Bob", "amount": 20}]
				}
			}`,
		},
		{
			name:         "No token",
			token:        "",
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"No token"}`,
		},
		{
			name:         "User not found",
			token:        "errorUser",
			expectedCode: fiber.StatusInternalServerError,
			expectedBody: `{"errors":"Failed to retrieve user info"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(AddClaimsToContextFromRequest())

			mockService := &MockUserInfoService{}

			app.Get("/api/info", func(c *fiber.Ctx) error {
				return GetInfo(c, mockService)
			})

			req := httptest.NewRequest("GET", "/api/info", nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error while making request: %v", err)
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var responseBody map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&responseBody)
			responseJSON, _ := json.Marshal(responseBody)

			assert.JSONEq(t, tt.expectedBody, string(responseJSON))
		})
	}
}
