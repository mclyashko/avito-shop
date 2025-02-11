package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/db"
)

var ErrNegativeSignTransaction = fmt.Errorf("negative sign")
var ErrInsufficientFunds = fmt.Errorf("insufficient funds")

type SendCoinsService interface {
	SendCoins(ctx context.Context, sender string, receiver string, amount int64) error
}

type SendCoinsServiceImpl struct {
	Service
	UserAccessor         db.UserAccessor
	CoinTransferAccessor db.CoinTransferAccessor
}

func (s *SendCoinsServiceImpl) SendCoins(ctx context.Context, sender string, receiver string, amount int64) error {
	if amount <= 0 {
		return ErrNegativeSignTransaction
	}

	return s.RunWithTx(ctx, func(tx pgx.Tx) error {
		senderUser, err := s.UserAccessor.GetUserByLoginTx(ctx, tx, sender)
		if err != nil {
			return err
		}

		_, err = s.UserAccessor.GetUserByLoginTx(ctx, tx, receiver)
		if err != nil {
			return err
		}

		if senderUser.Balance < amount {
			return ErrInsufficientFunds
		}

		if err := s.UserAccessor.UpdateUserBalanceTx(ctx, tx, sender, -amount); err != nil {
			return err
		}

		if err := s.UserAccessor.UpdateUserBalanceTx(ctx, tx, receiver, amount); err != nil {
			return err
		}

		if err := s.CoinTransferAccessor.InsertCoinTransferTx(ctx, tx, sender, receiver, amount); err != nil {
			return err
		}

		return nil
	})
}
