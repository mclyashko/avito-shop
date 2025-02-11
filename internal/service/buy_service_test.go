package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/model"
)

type mockService struct {
	Service
}

func (s *mockService) RunWithTx(ctx context.Context, txFunc func(tx pgx.Tx) error) error {
	return txFunc(nil)
}

type mockUserAccessor struct {
	users map[string]*model.User
	db.UserAccessor
}

func (m *mockUserAccessor) GetUserByLoginTx(ctx context.Context, tx pgx.Tx, username string) (*model.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *mockUserAccessor) UpdateUserBalanceTx(ctx context.Context, tx pgx.Tx, username string, amount int64) error {
	user, exists := m.users[username]
	if !exists {
		return fmt.Errorf("user not found")
	}
	user.Balance += amount
	return nil
}

type mockItemAccessor struct {
	items map[string]*model.Item
	db.ItemAccessor
}

func (m *mockItemAccessor) GetItemByNameTx(ctx context.Context, tx pgx.Tx, itemName string) (*model.Item, error) {
	item, exists := m.items[itemName]
	if !exists {
		return nil, fmt.Errorf("item not found")
	}
	return item, nil
}

type mockUserItemAccessor struct {
	userItems map[string][]string
	db.UserItemAccessor
}

func (m *mockUserItemAccessor) InsertUserItemTx(ctx context.Context, tx pgx.Tx, username string, itemName string) error {
	if _, exists := m.userItems[username]; !exists {
		m.userItems[username] = []string{}
	}
	m.userItems[username] = append(m.userItems[username], itemName)
	return nil
}

func TestBuyServiceImpl_BuyItem(t *testing.T) {
	mockService := &mockService{}
	mockUserAccessor := &mockUserAccessor{
		users: map[string]*model.User{
			"user1": {Login: "user1", Balance: 100},
		},
	}
	mockItemAccessor := &mockItemAccessor{
		items: map[string]*model.Item{
			"item1": {Name: "item1", Price: 50},
			"item2": {Name: "item2", Price: 200},
		},
	}
	mockUserItemAccessor := &mockUserItemAccessor{
		userItems: make(map[string][]string),
	}

	buyService := &BuyServiceImpl{
		Service:          mockService,
		UserAccessor:     mockUserAccessor,
		ItemAccessor:     mockItemAccessor,
		UserItemAccessor: mockUserItemAccessor,
	}

	tests := []struct {
		name      string
		username  string
		itemName  string
		wantErr   bool
		afterTest func()
	}{
		{
			name:     "successful purchase",
			username: "user1",
			itemName: "item1",
			wantErr:  false,
			afterTest: func() {
				user := mockUserAccessor.users["user1"]
				if user.Balance != 50 {
					t.Errorf("expected balance to be 50, got %d", user.Balance)
				}
				if len(mockUserItemAccessor.userItems["user1"]) != 1 {
					t.Errorf("expected 1 item in user's inventory, got %d", len(mockUserItemAccessor.userItems["user1"]))
				}
			},
		},
		{
			name:     "insufficient funds",
			username: "user1",
			itemName: "item2",
			wantErr:  true,
			afterTest: func() {
				user := mockUserAccessor.users["user1"]
				if user.Balance != 50 {
					t.Errorf("expected balance to remain 50, got %d", user.Balance)
				}
			},
		},
		{
			name:     "item not found",
			username: "user1",
			itemName: "item3",
			wantErr:  true,
			afterTest: func() {
				user := mockUserAccessor.users["user1"]
				if user.Balance != 50 {
					t.Errorf("expected balance to remain 50, got %d", user.Balance)
				}
			},
		},
		{
			name: "user not found",
			username: "aboba",
			itemName: "item1",
			wantErr: true,
			afterTest: func() {
				user := mockUserAccessor.users["user1"]
				if user.Balance != 50 {
					t.Errorf("expected balance to remain 50, got %d", user.Balance)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := buyService.BuyItem(context.Background(), tt.username, tt.itemName)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuyServiceImpl.BuyItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.afterTest()
		})
	}
}
