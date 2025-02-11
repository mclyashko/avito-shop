package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/model"
)

type mockCoinTransferAccessor struct {
	transfers []model.CoinTransfer
	db.CoinTransferAccessor
}

func (m *mockCoinTransferAccessor) InsertCoinTransferTx(ctx context.Context, tx pgx.Tx, sender string, receiver string, amount int64) error {
	m.transfers = append(m.transfers, model.CoinTransfer{SenderLogin: sender, ReceiverLogin: receiver, Amount: amount})
	return nil
}

func TestSendCoinsServiceImpl_SendCoins(t *testing.T) {
	mockService := &mockService{}
	mockUserAccessor := &mockUserAccessor{
		users: map[string]*model.User{
			"user1": {Login: "user1", Balance: 100},
			"user2": {Login: "user2", Balance: 50},
		},
	}
	mockCoinTransferAccessor := &mockCoinTransferAccessor{}

	sendCoinsService := &SendCoinsServiceImpl{
		Service:              mockService,
		UserAccessor:         mockUserAccessor,
		CoinTransferAccessor: mockCoinTransferAccessor,
	}

	tests := []struct {
		name      string
		sender    string
		receiver  string
		amount    int64
		wantErr   bool
		afterTest func()
	}{
		{
			name:     "successful transfer",
			sender:   "user1",
			receiver: "user2",
			amount:   50,
			wantErr:  false,
			afterTest: func() {
				if mockUserAccessor.users["user1"].Balance != 50 {
					t.Errorf("expected user1 balance to be 50, got %d", mockUserAccessor.users["user1"].Balance)
				}
				if mockUserAccessor.users["user2"].Balance != 100 {
					t.Errorf("expected user2 balance to be 100, got %d", mockUserAccessor.users["user2"].Balance)
				}
				if len(mockCoinTransferAccessor.transfers) != 1 {
					t.Errorf("expected 1 coin transfer, got %d", len(mockCoinTransferAccessor.transfers))
				}
			},
		},
		{
			name:     "insufficient funds",
			sender:   "user1",
			receiver: "user2",
			amount:   150,
			wantErr:  true,
			afterTest: func() {
				if mockUserAccessor.users["user1"].Balance != 50 {
					t.Errorf("expected user1 balance to remain 50, got %d", mockUserAccessor.users["user1"].Balance)
				}
			},
		},
		{
			name:     "negative amount",
			sender:   "user1",
			receiver: "user2",
			amount:   -50,
			wantErr:  true,
			afterTest: func() {
				if mockUserAccessor.users["user1"].Balance != 50 {
					t.Errorf("expected user1 balance to remain 50, got %d", mockUserAccessor.users["user1"].Balance)
				}
			},
		},
		{
			name:     "sender not found",
			sender:   "user3",
			receiver: "user2",
			amount:   50,
			wantErr:  true,
			afterTest: func() {
				if len(mockCoinTransferAccessor.transfers) != 1 {
					t.Errorf("expected 1 transfer, got %d", len(mockCoinTransferAccessor.transfers))
				}
			},
		},
		{
			name:     "receiver not found",
			sender:   "user1",
			receiver: "user3",
			amount:   50,
			wantErr:  true,
			afterTest: func() {
				if len(mockCoinTransferAccessor.transfers) != 1 {
					t.Errorf("expected 1 transfer, got %d", len(mockCoinTransferAccessor.transfers))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sendCoinsService.SendCoins(context.Background(), tt.sender, tt.receiver, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("SendCoinsServiceImpl.SendCoins() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.afterTest()
		})
	}
}
