package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/mclyashko/avito-shop/internal/model"
)

func (m *mockUserAccessor) GetUserByLogin(ctx context.Context, username string) (*model.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *mockUserItemAccessor) GetUserItemsByUsername(ctx context.Context, username string) ([]model.UserItem, error) {
	userItems, exists := m.userItems[username]
	if !exists {
		return nil, fmt.Errorf("user items not found")
	}

	result := make([]model.UserItem, len(userItems))
	for i, item := range userItems {
		result[i] = model.UserItem{
			ItemName: item,
		}
	}
	return result, nil
}

func (m *mockCoinTransferAccessor) GetUserTransactionHistory(ctx context.Context, username string) ([]model.CoinTransfer, []model.CoinTransfer, error) {
	var receivedTransfers []model.CoinTransfer
	var sentTransfers []model.CoinTransfer
	for _, transfer := range m.transfers {
		if transfer.ReceiverLogin == username {
			receivedTransfers = append(receivedTransfers, transfer)
		}
		if transfer.SenderLogin == username {
			sentTransfers = append(sentTransfers, transfer)
		}
	}
	return receivedTransfers, sentTransfers, nil
}

func TestUserInfoServiceImp_GetUserInfo(t *testing.T) {
	mockUserAccessor := &mockUserAccessor{
		users: map[string]*model.User{
			"user1": {Login: "user1", Balance: 100},
		},
	}
	mockUserItemAccessor := &mockUserItemAccessor{
		userItems: map[string][]string{
			"user1": {"item1"},
		},
	}
	mockCoinTransferAccessor := &mockCoinTransferAccessor{
		transfers: []model.CoinTransfer{
			{
				SenderLogin:   "user1",
				ReceiverLogin: "user2",
			},
			{
				SenderLogin:   "user1",
				ReceiverLogin: "user3",
			},
			{
				SenderLogin:   "user3",
				ReceiverLogin: "user1",
			},
		},
	}

	userInfoService := &UserInfoServiceImp{
		UserAccessor:         mockUserAccessor,
		UserItemAccessor:     mockUserItemAccessor,
		CoinTransferAccessor: mockCoinTransferAccessor,
	}

	tests := []struct {
		name         string
		username     string
		expectedErr  bool
		expectedData struct {
			balance           int64
			userItems         int
			receivedTransfers int
			sentTransfers     int
		}
	}{
		{
			name:        "successful user info retrieval",
			username:    "user1",
			expectedErr: false,
			expectedData: struct {
				balance           int64
				userItems         int
				receivedTransfers int
				sentTransfers     int
			}{
				balance:           100,
				userItems:         1,
				receivedTransfers: 1,
				sentTransfers:     2,
			},
		},
		{
			name:        "user not found",
			username:    "user2",
			expectedErr: true,
			expectedData: struct {
				balance           int64
				userItems         int
				receivedTransfers int
				sentTransfers     int
			}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, userItems, receivedTransfers, sentTransfers, err := userInfoService.GetUserInfo(context.Background(), tt.username)

			if (err != nil) != tt.expectedErr {
				t.Errorf("UserInfoServiceImp.GetUserInfo() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}

			if !tt.expectedErr {
				if *balance != tt.expectedData.balance {
					t.Errorf("expected balance %d, got %d", tt.expectedData.balance, *balance)
				}
				if len(userItems) != tt.expectedData.userItems {
					t.Errorf("expected %d user items, got %d", tt.expectedData.userItems, len(userItems))
				}
				if len(receivedTransfers) != tt.expectedData.receivedTransfers {
					t.Errorf("expected %d received transfers, got %d", tt.expectedData.receivedTransfers, len(receivedTransfers))
				}
				if len(sentTransfers) != tt.expectedData.sentTransfers {
					t.Errorf("expected %d sent transfers, got %d", tt.expectedData.sentTransfers, len(sentTransfers))
				}
			}
		})
	}
}
