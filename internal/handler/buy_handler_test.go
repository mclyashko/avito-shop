package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

func AddClaimsToContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token != "" {
			c.Locals("claims", &service.JWTClaims{
				Username: "testUser",
			})
		}
		return c.Next()
	}
}

type MockBuyService struct{}

func (s *MockBuyService) BuyItem(ctx context.Context, username, itemName string) error {
	if itemName == "invalidItem" {
		return db.ErrItemNotFound
	}
	if itemName == "insufficientFunds" {
		return service.ErrInsufficientFunds
	}
	return nil
}

func TestBuyItemHandler(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		itemName     string
		mockResponse error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful purchase",
			token:        "validToken",
			itemName:     "validItem",
			mockResponse: nil,
			expectedCode: fiber.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Item not found",
			token:        "validToken",
			itemName:     "invalidItem",
			mockResponse: db.ErrItemNotFound,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Item not found"}`,
		},
		{
			name:         "Insufficient funds",
			token:        "validToken",
			itemName:     "insufficientFunds",
			mockResponse: service.ErrInsufficientFunds,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"Insufficient funds"}`,
		},
		{
			name:         "Missing token",
			token:        "",
			itemName:     "validItem",
			mockResponse: nil,
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"errors":"No token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(AddClaimsToContext())

			mockBuyService := &MockBuyService{}

			app.Post("/api/buy/:item", func(c *fiber.Ctx) error {
				return BuyItemHandler(c, mockBuyService)
			})

			body := []byte(`{}`)
			req := newTestRequest("POST", "/api/buy/"+tt.itemName, body)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error while making the request: %v", err)
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedBody == "" {
				assert.Empty(t, resp.Body)
			} else {
				var responseBody map[string]interface{}
				_ = json.NewDecoder(resp.Body).Decode(&responseBody)
				responseJSON, _ := json.Marshal(responseBody)

				assert.JSONEq(t, tt.expectedBody, string(responseJSON))
			}
		})
	}
}
